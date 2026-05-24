#!/usr/bin/env python3
import os
import sys
import shutil
import tempfile
import subprocess
import argparse
import requests

# GitHub repositories
PAPERCLIP_UPSTREAM = "paperclipai/paperclip"

# Default Portainer configuration
DEFAULT_PORTAINER_URL = "https://painel.wasend.com.br"
DEFAULT_ENDPOINT_ID = 1

def parse_version(ver_str):
    """
    Parse a version string (e.g. 'v2026.428.0', 'v2026.429.0-fixed') into a tuple of integers.
    """
    # Strip prefixes and suffixes
    ver_str = ver_str.lstrip('v').split('-')[0].strip()
    parts = []
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

def get_running_paperclip_version():
    """
    Connect to Portainer and find the running version (tag) of heltonfraga/paperclipai.
    If unavailable, returns None.
    """
    url = os.environ.get("PORTAINER_URL", DEFAULT_PORTAINER_URL)
    username = os.environ.get("PORTAINER_USERNAME", "CortexxRick")
    password = os.environ.get("PORTAINER_PASSWORD", "14850d130b7308@eadB")
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

        print(f"[*] Fetching all Docker services on endpoint {endpoint_id}...")
        headers = {"Authorization": f"Bearer {token}"}
        services_resp = requests.get(f"{url}/api/endpoints/{endpoint_id}/docker/services", headers=headers, timeout=15)
        services_resp.raise_for_status()
        services = services_resp.json()

        for service in services:
            image = service.get("Spec", {}).get("TaskTemplate", {}).get("ContainerSpec", {}).get("Image", "")
            if "heltonfraga/paperclipai" in image:
                print(f"[*] Found running Paperclip service: {service.get('Spec', {}).get('Name')} using image {image}")
                # Strip @sha256 digest first if present
                image_no_digest = image.split("@")[0]
                if ":" in image_no_digest:
                    return image_no_digest.split(":")[-1]
    except Exception as e:
        print(f"[!] Failed to get running Paperclip version from Portainer: {e}")
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

def check_existing_issue(repo, title_query):
    """
    Check if there is already an open issue with the given title query in the repo.
    """
    url = f"https://api.github.com/repos/{repo}/issues?state=open"
    try:
        issues = github_api_request("GET", url)
        if issues:
            for issue in issues:
                if title_query in issue.get("title", ""):
                    return issue.get("html_url")
    except Exception as e:
        print(f"[!] Error checking existing issues: {e}")
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

def run_cmd(args, cwd=None, check=True):
    """
    Run a terminal command and return stdout.
    """
    print(f"[*] Executing command: {' '.join(args)} (cwd={cwd})")
    result = subprocess.run(args, capture_output=True, text=True, cwd=cwd)
    if check and result.returncode != 0:
        print(f"[!] Command failed: {' '.join(args)}")
        print(f"Stderr: {result.stderr}")
        raise subprocess.CalledProcessError(result.returncode, args, result.stdout, result.stderr)
    return result.stdout.strip()

def apply_code_patches(directory):
    """
    Recursively walk the cloned repo and replace hardcoded Anthropic/OpenAI API URLs
    with environment variables.
    """
    print(f"[*] Applying environment patches recursively to {directory}...")
    replacements = {
        '"https://api.anthropic.com"': 'process.env.ANTHROPIC_BASE_URL || "https://api.anthropic.com"',
        "'https://api.anthropic.com'": "process.env.ANTHROPIC_BASE_URL || 'https://api.anthropic.com'",
        '"https://api.openai.com"': 'process.env.OPENAI_BASE_URL || "https://api.openai.com"',
        "'https://api.openai.com'": "process.env.OPENAI_BASE_URL || 'https://api.openai.com'",
        '"https://api.anthropic.com/v1/messages"': '(process.env.ANTHROPIC_BASE_URL || "https://api.anthropic.com") + "/v1/messages"',
        "'https://api.anthropic.com/v1/messages'": "(process.env.ANTHROPIC_BASE_URL || 'https://api.anthropic.com') + '/v1/messages'",
        '"https://api.anthropic.com/api/oauth/usage"': '(process.env.ANTHROPIC_BASE_URL || "https://api.anthropic.com") + "/api/oauth/usage"',
        "'https://api.anthropic.com/api/oauth/usage'": "(process.env.ANTHROPIC_BASE_URL || 'https://api.anthropic.com') + '/api/oauth/usage'",
        '"https://api.openai.com/v1/models"': '(process.env.OPENAI_BASE_URL || "https://api.openai.com") + "/v1/models"',
        "'https://api.openai.com/v1/models'": "(process.env.OPENAI_BASE_URL || 'https://api.openai.com') + '/v1/models'",
        '"https://chatgpt.com/backend-api/wham/usage"': '(process.env.CHATGPT_WHAM_USAGE_URL || (process.env.OPENAI_BASE_URL ? (process.env.OPENAI_BASE_URL.endsWith("/v1") ? process.env.OPENAI_BASE_URL.slice(0, -3) : process.env.OPENAI_BASE_URL) + "/backend-api/wham/usage" : "https://chatgpt.com/backend-api/wham/usage"))',
        "'https://chatgpt.com/backend-api/wham/usage'": "(process.env.CHATGPT_WHAM_USAGE_URL || (process.env.OPENAI_BASE_URL ? (process.env.OPENAI_BASE_URL.endsWith('/v1') ? process.env.OPENAI_BASE_URL.slice(0, -3) : process.env.OPENAI_BASE_URL) + '/backend-api/wham/usage' : 'https://chatgpt.com/backend-api/wham/usage'))",
    }


    count = 0
    for root, dirs, files in os.walk(directory):
        # Skip node_modules, git, and other common build/system dirs
        if any(ignored in root for ignored in [".git", "node_modules", "dist", "build"]):
            continue

        for filename in files:
            # Only patch code files (js, ts, jsx, tsx, json, py)
            if not filename.endswith(('.js', '.ts', '.jsx', '.tsx', '.json', '.py')):
                continue

            filepath = os.path.join(root, filename)
            try:
                with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
                    content = f.read()

                modified = False
                for target, replacement in replacements.items():
                    if target in content:
                        content = content.replace(target, replacement)
                        modified = True

                if modified:
                    with open(filepath, 'w', encoding='utf-8') as f:
                        f.write(content)
                    print(f"  [+] Patched file: {os.path.relpath(filepath, directory)}")
                    count += 1
            except Exception as e:
                print(f"  [!] Failed to patch file {filepath}: {e}")

    print(f"[+] Patch application finished. Modified {count} files.")
    return count

def run_docker_build_and_push(temp_dir, version_tag, dry_run=False):
    """
    Build multi-arch docker image and push to Docker Hub.
    """
    target_image_fixed = f"heltonfraga/paperclipai:{version_tag}-fixed"
    target_image_latest = "heltonfraga/paperclipai:latest"

    if dry_run:
        print(f"[Dry-run] Would build and push images:")
        print(f"  - {target_image_fixed}")
        print(f"  - {target_image_latest}")
        return True

    # 1. Login to Docker Hub if credentials are provided in env
    dh_username = os.environ.get("DOCKERHUB_USERNAME")
    dh_token = os.environ.get("DOCKERHUB_TOKEN")

    if dh_username and dh_token:
        print("[*] Logging in to Docker Hub...")
        login_process = subprocess.run(
            ["docker", "login", "-u", dh_username, "--password-stdin"],
            input=dh_token,
            text=True,
            capture_output=True
        )
        if login_process.returncode != 0:
            print(f"[!] Docker Hub login failed: {login_process.stderr}")
            return False
        print("[+] Logged in successfully.")
    else:
        print("[*] Docker Hub credentials not fully provided in environment, skipping login step.")

    # 2. Build multi-arch image and push
    print(f"[*] Building and pushing Docker image using Buildx for {version_tag}...")
    try:
        # Check if buildx builder exists and use it
        builders = run_cmd(["docker", "buildx", "ls"])
        if "multiarch-builder" not in builders:
            print("[*] Creating docker buildx multi-arch builder...")
            run_cmd(["docker", "buildx", "create", "--name", "multiarch-builder", "--use"])
            run_cmd(["docker", "buildx", "inspect", "--bootstrap"])
        else:
            run_cmd(["docker", "buildx", "use", "multiarch-builder"])

        # Execute build and push
        build_args = [
            "docker", "buildx", "build",
            "--platform", "linux/amd64,linux/arm64",
            "-t", target_image_fixed,
            "-t", target_image_latest,
            "--push",
            "."
        ]
        run_cmd(build_args, cwd=temp_dir)
        print(f"[+] Build and push succeeded for {target_image_fixed}")
        return True
    except Exception as e:
        print(f"[!] Docker buildx build failed: {e}")
        return False

def process_paperclip_sync(latest_tag, release_notes, repo_fullname, dry_run=False):
    """
    Clone, patch, build, push, and create issue.
    """
    title_query = f"Build Paperclip {latest_tag}-fixed"
    existing = check_existing_issue(repo_fullname, title_query)
    if existing:
        print(f"[*] Patched release build issue already exists: {existing}")
        return

    print(f"[*] Starting clone/patch/build pipeline for Paperclip {latest_tag}...")

    # Create temporary directory for cloning
    temp_dir = tempfile.mkdtemp()
    try:
        # 1. Clone upstream
        print(f"[*] Cloning {PAPERCLIP_UPSTREAM} at tag {latest_tag} to {temp_dir}...")
        run_cmd(["git", "clone", "--depth", "1", "--branch", latest_tag, f"https://github.com/{PAPERCLIP_UPSTREAM}.git", temp_dir])

        # 2. Apply patch
        apply_code_patches(temp_dir)

        # 3. Build & Push Docker image
        success = run_docker_build_and_push(temp_dir, latest_tag, dry_run)

        # 4. Create Issue notifying user
        if success and not dry_run:
            issue_title = f"Build Paperclip {latest_tag}-fixed"
            issue_body = (
                f"A new official release of Paperclip has been detected: **{latest_tag}**.\n\n"
                f"The patched image has been successfully compiled and pushed to Docker Hub:\n"
                f"- `heltonfraga/paperclipai:{latest_tag}-fixed`\n"
                f"- `heltonfraga/paperclipai:latest`\n\n"
                f"### What to do:\n"
                f"1. Open Portainer.\n"
                f"2. Go to the **paperclip** stack.\n"
                f"3. Update the image to `heltonfraga/paperclipai:{latest_tag}-fixed`.\n"
                f"4. Click **Update the stack**.\n\n"
                f"### Upstream Release Notes:\n\n{release_notes}"
            )
            create_github_issue(repo_fullname, issue_title, issue_body)

    finally:
        # Clean up temp folder
        print(f"[*] Cleaning up temporary directory {temp_dir}...")
        shutil.rmtree(temp_dir)

def main():
    parser = argparse.ArgumentParser(description="Monitor, patch, and build Paperclip updates")
    parser.add_argument("--dry-run", action="store_true", help="Run without pushing Docker image or creating Issues")
    parser.add_argument("--repo", default="HeltonFraga01/nexusmind", help="Our fork/main repository fullname on GitHub")
    args = parser.parse_args()

    dry_run = args.dry_run
    repo_fullname = args.repo

    print(f"=== Starting Paperclip Upstream Monitor & Patch Automation (Dry-run={dry_run}) ===")

    # 1. Get running version from Portainer
    running_version = get_running_paperclip_version()
    if not running_version:
        running_version = "v2026.428.0" # fallback based on stack YAML
        print(f"[*] Portainer check failed/skipped. Using fallback running version: {running_version}")
    else:
        print(f"[*] Running Paperclip version in Portainer: {running_version}")

    # 2. Get latest upstream release
    try:
        latest_tag, release_notes, release_url = get_latest_upstream_release(PAPERCLIP_UPSTREAM)
        print(f"[*] Latest upstream Paperclip version: {latest_tag} ({release_url})")

        # Compare versions
        if parse_version(latest_tag) > parse_version(running_version):
            print(f"[!] Version drift detected! Upstream {latest_tag} > Running {running_version}")
            process_paperclip_sync(latest_tag, release_notes, repo_fullname, dry_run)
        else:
            print("[+] Paperclip is up to date with upstream.")
    except Exception as e:
        print(f"[!] Error checking Paperclip upstream: {e}")

    print("\n=== Paperclip Upstream Monitor Completed ===")

if __name__ == "__main__":
    main()
