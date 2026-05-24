#!/usr/bin/env python3
import os
import sys
import subprocess
import argparse
import requests

# GitHub repositories
SUB2API_UPSTREAM = "Wei-Shaw/sub2api"
OMNIROUTE_UPSTREAM = "diegosouzapw/OmniRoute"

# Default Portainer configuration for OmniRoute query
DEFAULT_PORTAINER_URL = "https://painel.wasend.com.br"
DEFAULT_SERVICE_ID = "fnnzg5nqoc1wustre216jkph1" # cortexx-omniroute_cortexx-omniroute
DEFAULT_ENDPOINT_ID = 1

def parse_version(ver_str):
    """
    Parse a version string (e.g. 'v0.1.130', '3.8.2') into a tuple of integers.
    """
    ver_str = ver_str.lstrip('v').strip()
    parts = []
    # Split by '.' and extract leading digits from each part
    for p in ver_str.split('.'):
        num = ""
        for char in p:
            if char.isdigit():
                num += char
            else:
                break
        if num:
            parts.append(int(num))
        else:
            parts.append(0)
    return tuple(parts)

def get_latest_upstream_release(repo):
    """
    Fetch the latest release tag name and notes from the GitHub API.
    """
    print(f"[*] Fetching latest release for upstream repo: {repo}...")
    url = f"https://api.github.com/repos/{repo}/releases/latest"
    headers = {}
    token = os.environ.get("GITHUB_TOKEN") or os.environ.get("GH_TOKEN")
    if token:
        headers["Authorization"] = f"token {token}"
    
    resp = requests.get(url, headers=headers, timeout=15)
    resp.raise_for_status()
    data = resp.json()
    return data.get("tag_name"), data.get("body", ""), data.get("html_url", "")

def get_local_sub2api_version(version_file_path):
    """
    Read the local version of Sub2API from the VERSION file.
    """
    if not os.path.exists(version_file_path):
        print(f"[!] Version file not found at {version_file_path}, defaulting to 0.0.0")
        return "0.0.0"
    with open(version_file_path, "r", encoding="utf-8") as f:
        return f.read().strip()

def get_running_omniroute_version():
    """
    Get the running OmniRoute version from Portainer.
    If Portainer is unavailable or credentials are not set, returns None.
    """
    url = os.environ.get("PORTAINER_URL", DEFAULT_PORTAINER_URL)
    username = os.environ.get("PORTAINER_USERNAME", "CortexxRick")
    password = os.environ.get("PORTAINER_PASSWORD", "14850d130b7308@eadB")
    service_id = os.environ.get("PORTAINER_SERVICE_ID", DEFAULT_SERVICE_ID)
    endpoint_id_str = os.environ.get("PORTAINER_ENDPOINT_ID")
    endpoint_id = int(endpoint_id_str) if endpoint_id_str else DEFAULT_ENDPOINT_ID

    if not username or not password:
        print("[*] Portainer credentials not fully set, skipping Portainer check.")
        return None

    try:
        print(f"[*] Authenticating with Portainer at {url}...")
        auth_resp = requests.post(f"{url}/api/auth", json={
            "username": username,
            "password": password
        }, timeout=10)
        auth_resp.raise_for_status()
        token = auth_resp.json()["jwt"]

        print(f"[*] Fetching config for OmniRoute service {service_id}...")
        headers = {"Authorization": f"Bearer {token}"}
        service_resp = requests.get(f"{url}/api/endpoints/{endpoint_id}/docker/services/{service_id}", headers=headers, timeout=10)
        service_resp.raise_for_status()
        service = service_resp.json()

        image = service.get("Spec", {}).get("TaskTemplate", {}).get("ContainerSpec", {}).get("Image", "")
        print(f"[*] Found image in Portainer spec: {image}")
        # Strip @sha256 digest first if present
        image_no_digest = image.split("@")[0]
        if ":" in image_no_digest:
            return image_no_digest.split(":")[-1]
    except Exception as e:
        print(f"[!] Failed to get OmniRoute version from Portainer: {e}")
    return None

def github_api_request(method, url, json_data=None):
    """
    Make an authenticated request to the GitHub API.
    """
    token = os.environ.get("GITHUB_TOKEN") or os.environ.get("GH_TOKEN")
    if not token:
        print("[!] No GITHUB_TOKEN or GH_TOKEN found, cannot make modification requests.")
        return None
    headers = {
        "Authorization": f"token {token}",
        "Accept": "application/vnd.github.v3+json"
    }
    resp = requests.request(method, url, headers=headers, json=json_data, timeout=15)
    if resp.status_code >= 400:
        print(f"[!] GitHub API error ({resp.status_code}): {resp.text}")
    resp.raise_for_status()
    return resp.json()

def check_existing_issue_or_pr(repo, title_query):
    """
    Check if there is already an open issue or PR with the given title query in the repo.
    """
    url = f"https://api.github.com/repos/{repo}/issues?state=open"
    try:
        issues = github_api_request("GET", url)
        if issues:
            for issue in issues:
                if title_query in issue.get("title", ""):
                    return issue.get("html_url")
    except Exception as e:
        print(f"[!] Error checking existing issues/PRs: {e}")
    return None

def create_github_issue(repo, title, body):
    """
    Create a GitHub Issue in the specified repository.
    """
    url = f"https://api.github.com/repos/{repo}/issues"
    payload = {
        "title": title,
        "body": body
    }
    print(f"[*] Creating GitHub Issue in {repo}: '{title}'...")
    try:
        data = github_api_request("POST", url, payload)
        if data:
            print(f"[+] Issue created successfully: {data.get('html_url')}")
            return data.get("html_url")
    except Exception as e:
        print(f"[!] Failed to create GitHub Issue: {e}")
    return None

def create_github_pr(repo, title, body, head, base="main"):
    """
    Create a GitHub Pull Request in the specified repository.
    """
    url = f"https://api.github.com/repos/{repo}/pulls"
    payload = {
        "title": title,
        "body": body,
        "head": head,
        "base": base
    }
    print(f"[*] Creating GitHub Pull Request in {repo}: '{title}' from branch '{head}' to '{base}'...")
    try:
        data = github_api_request("POST", url, payload)
        if data:
            print(f"[+] Pull Request created successfully: {data.get('html_url')}")
            return data.get("html_url")
    except Exception as e:
        print(f"[!] Failed to create Pull Request: {e}")
    return None

def run_cmd(args, check=True):
    """
    Run a terminal command and return stdout.
    """
    print(f"[*] Executing command: {' '.join(args)}")
    result = subprocess.run(args, capture_output=True, text=True)
    if check and result.returncode != 0:
        print(f"[!] Command failed: {' '.join(args)}")
        print(f"Stderr: {result.stderr}")
        raise subprocess.CalledProcessError(result.returncode, args, result.stdout, result.stderr)
    return result.stdout.strip()

def attempt_sub2api_merge(target_tag, release_notes, repo_fullname, dry_run=False):
    """
    Attempt to fetch the upstream tag, create a sync branch, merge, and push/open PR.
    """
    title_query = f"Sync upstream Sub2API {target_tag}"
    existing = check_existing_issue_or_pr(repo_fullname, title_query)
    if existing:
        print(f"[*] Sync PR or Issue already exists: {existing}")
        return

    if dry_run:
        print(f"[Dry-run] Would attempt to merge upstream tag {target_tag} into branch sync/upstream-{target_tag}")
        return

    branch_name = f"sync/upstream-{target_tag}"
    print(f"[*] Starting upstream sync process for tag {target_tag}...")
    
    try:
        # Configure git user if running in CI
        if os.environ.get("GITHUB_ACTIONS") == "true":
            run_cmd(["git", "config", "user.name", "github-actions[bot]"])
            run_cmd(["git", "config", "user.email", "41898282+github-actions[bot]@users.noreply.github.com"])

        # Add upstream remote if it doesn't exist
        remotes = run_cmd(["git", "remote"])
        if "upstream" not in remotes.split():
            run_cmd(["git", "remote", "add", "upstream", f"https://github.com/{SUB2API_UPSTREAM}.git"])

        # Fetch upstream and tags
        run_cmd(["git", "fetch", "upstream", "main", "--tags"])

        # Ensure we are starting from main
        run_cmd(["git", "checkout", "main"])
        run_cmd(["git", "pull", "origin", "main"])

        # Create and checkout sync branch
        run_cmd(["git", "checkout", "-b", branch_name])

        # Attempt to merge the tag
        print(f"[*] Attempting to merge tag {target_tag} into {branch_name}...")
        merge_proc = subprocess.run(["git", "merge", target_tag, "-m", f"Merge upstream tag {target_tag}"], capture_output=True, text=True)
        
        if merge_proc.returncode == 0:
            print("[+] Merge completed with NO conflicts!")
            # Push the branch to origin
            run_cmd(["git", "push", "origin", branch_name, "--force"])
            
            # Create Pull Request
            pr_body = (
                f"Automatic sync with upstream `{SUB2API_UPSTREAM}` release `{target_tag}`.\n\n"
                f"### Upstream Release Notes:\n\n{release_notes}\n\n"
                f"Please review and merge if everything works."
            )
            create_github_pr(repo_fullname, f"Sync upstream Sub2API {target_tag}", pr_body, branch_name)
        else:
            print("[!] Merge failed with conflicts! Aborting merge and reporting...")
            run_cmd(["git", "merge", "--abort"])
            
            # Create an issue about the conflict
            issue_title = f"[Conflict] Sync upstream Sub2API {target_tag} failed"
            issue_body = (
                f"Automatic merge of upstream tag `{target_tag}` failed due to conflicts.\n\n"
                f"Please resolve conflicts manually:\n"
                f"```bash\n"
                f"git checkout main\n"
                f"git pull origin main\n"
                f"git fetch upstream --tags\n"
                f"git merge {target_tag}\n"
                f"# resolve conflicts, commit, and push\n"
                f"```\n\n"
                f"### Conflict Output:\n```\n{merge_proc.stderr or merge_proc.stdout}\n```"
            )
            create_github_issue(repo_fullname, issue_title, issue_body)

        # Return to main branch
        run_cmd(["git", "checkout", "main"])

    except Exception as e:
        print(f"[!] Error during git merge operations: {e}")
        # Make sure we clean up and return to main
        try:
            run_cmd(["git", "merge", "--abort"], check=False)
            run_cmd(["git", "checkout", "main"], check=False)
        except:
            pass

def main():
    parser = argparse.ArgumentParser(description="Monitor upstream updates for Sub2API and OmniRoute")
    parser.add_argument("--dry-run", action="store_true", help="Print findings without making changes on GitHub")
    parser.add_argument("--repo", default="HeltonFraga01/nexusmind", help="Our fork repository fullname on GitHub")
    parser.add_argument("--version-file", default="backend/cmd/server/VERSION", help="Path to Sub2API local VERSION file")
    args = parser.parse_args()

    dry_run = args.dry_run
    repo_fullname = args.repo
    version_file = args.version_file

    print(f"=== Starting Upstream Update Monitor (Dry-run={dry_run}) ===")

    # ----------------------------------------------------
    # 1. Monitor Sub2API
    # ----------------------------------------------------
    print("\n--- Checking Sub2API ---")
    local_sub2api = get_local_sub2api_version(version_file)
    print(f"[*] Local Sub2API version: {local_sub2api}")

    try:
        latest_tag, sub2api_notes, sub2api_url = get_latest_upstream_release(SUB2API_UPSTREAM)
        clean_tag = latest_tag.lstrip('v').strip()
        print(f"[*] Latest upstream Sub2API version: {clean_tag} ({sub2api_url})")

        if parse_version(clean_tag) > parse_version(local_sub2api):
            print(f"[!] Version drift detected! Upstream {clean_tag} > Local {local_sub2api}")
            attempt_sub2api_merge(latest_tag, sub2api_notes, repo_fullname, dry_run)
        else:
            print("[+] Sub2API is up to date with upstream.")
    except Exception as e:
        print(f"[!] Error checking Sub2API upstream: {e}")

    # ----------------------------------------------------
    # 2. Monitor OmniRoute
    # ----------------------------------------------------
    print("\n--- Checking OmniRoute ---")
    running_omniroute = get_running_omniroute_version()
    if not running_omniroute:
        # Fallback to last known deployed version if Portainer check failed/skipped
        running_omniroute = "3.8.2"
        print(f"[*] Portainer query skipped or failed. Using fallback OmniRoute version: {running_omniroute}")
    else:
        print(f"[*] Running OmniRoute version in Portainer Swarm: {running_omniroute}")

    try:
        latest_omni_tag, omni_notes, omni_url = get_latest_upstream_release(OMNIROUTE_UPSTREAM)
        clean_omni_tag = latest_omni_tag.lstrip('v').strip()
        print(f"[*] Latest upstream OmniRoute version: {clean_omni_tag} ({omni_url})")

        if parse_version(clean_omni_tag) > parse_version(running_omniroute):
            print(f"[!] Version drift detected! Upstream {clean_omni_tag} > Running {running_omniroute}")
            
            title_query = f"Update OmniRoute to {latest_omni_tag}"
            existing = check_existing_issue_or_pr(repo_fullname, title_query)
            if existing:
                print(f"[*] Update issue already exists: {existing}")
            else:
                issue_title = f"Update OmniRoute to {latest_omni_tag}"
                issue_body = (
                    f"A new version of OmniRoute is available upstream: **{latest_omni_tag}** (currently running: **{running_omniroute}**).\n\n"
                    f"Please run the Portainer deployment script to update:\n"
                    f"```bash\n"
                    f"python3 scripts/update_omniroute.py --version {clean_omni_tag}\n"
                    f"```\n\n"
                    f"### Upstream Release Notes:\n\n{omni_notes}\n\n"
                    f"Upstream release link: {omni_url}"
                )
                if dry_run:
                    print(f"[Dry-run] Would create issue in {repo_fullname}: '{issue_title}'")
                else:
                    create_github_issue(repo_fullname, issue_title, issue_body)
        else:
            print("[+] OmniRoute is up to date with upstream.")
    except Exception as e:
        print(f"[!] Error checking OmniRoute upstream: {e}")

    print("\n=== Upstream Update Monitor Completed ===")

if __name__ == "__main__":
    main()
