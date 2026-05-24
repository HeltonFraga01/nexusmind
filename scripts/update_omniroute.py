#!/usr/bin/env python3
import argparse
import json
import os
import sys
import time
import requests

# Default settings
DEFAULT_PORTAINER_URL = "https://painel.wasend.com.br"
DEFAULT_SERVICE_ID = "fnnzg5nqoc1wustre216jkph1" # cortexx-omniroute_cortexx-omniroute
DEFAULT_ENDPOINT_ID = 1

def get_env_or_default(key, default):
    return os.environ.get(key, default)

def authenticate(url, username, password):
    print(f"[*] Authenticating with Portainer at {url}...")
    resp = requests.post(f"{url}/api/auth", json={
        "username": username,
        "password": password
    }, timeout=15)
    resp.raise_for_status()
    return resp.json()["jwt"]

def get_service(url, token, service_id, endpoint_id):
    headers = {"Authorization": f"Bearer {token}"}
    resp = requests.get(f"{url}/api/endpoints/{endpoint_id}/docker/services/{service_id}", headers=headers, timeout=15)
    resp.raise_for_status()
    return resp.json()

def update_service(url, token, service_id, endpoint_id, version_index, spec):
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }
    update_url = f"{url}/api/endpoints/{endpoint_id}/docker/services/{service_id}/update?version={version_index}"
    resp = requests.post(update_url, headers=headers, json=spec, timeout=30)
    resp.raise_for_status()
    return resp.json()

def get_containers(url, token, endpoint_id):
    headers = {"Authorization": f"Bearer {token}"}
    resp = requests.get(f"{url}/api/endpoints/{endpoint_id}/docker/containers/json?all=1", headers=headers, timeout=15)
    resp.raise_for_status()
    return resp.json()

def get_container_logs(url, token, container_id, endpoint_id):
    headers = {"Authorization": f"Bearer {token}"}
    log_url = f"{url}/api/endpoints/{endpoint_id}/docker/containers/{container_id}/logs?stdout=1&stderr=1&tail=50"
    resp = requests.get(log_url, headers=headers, timeout=15)
    resp.raise_for_status()
    return resp.text

def monitor_deployment(url, token, endpoint_id, old_container_id, timeout_sec=180):
    print("[*] Monitoring deployment status...")
    start_time = time.time()
    while time.time() - start_time < timeout_sec:
        try:
            containers = get_containers(url, token, endpoint_id)
        except Exception as e:
            print(f"[!] Error listing containers: {e}")
            time.sleep(5)
            continue
            
        new_container = None
        for c in containers:
            names = c.get("Names", [])
            name = names[0] if names else ""
            if "cortexx-omniroute_cortexx-omniroute" in name and c.get("Id") != old_container_id:
                new_container = c
                break
                
        if new_container:
            state = new_container.get("State")
            status = new_container.get("Status")
            print(f"[*] New container found: State={state}, Status={status}")
            if state == "running":
                if "healthy" in status:
                    print("[+] Deployment completed successfully and container is healthy!")
                    return True
                elif "unhealthy" in status:
                    print("[!] Container is marked unhealthy! Fetching logs...")
                    logs = get_container_logs(url, token, new_container["Id"], endpoint_id)
                    print("\n--- CONTAINER LOGS ---")
                    print(logs)
                    return False
            elif state == "exited":
                print(f"[!] Container exited with state {state}! Fetching logs...")
                logs = get_container_logs(url, token, new_container["Id"], endpoint_id)
                print("\n--- CONTAINER LOGS ---")
                print(logs)
                return False
        else:
            print("[*] Waiting for new container to spawn...")
            
        time.sleep(5)
    print("[!] Timeout waiting for deployment to complete.")
    return False

def main():
    parser = argparse.ArgumentParser(description="Update cortexx-omniroute service in Portainer")
    parser.add_argument("--version", required=True, help="Target OmniRoute version to deploy (e.g. 3.8.2)")
    parser.add_argument("--url", default=None, help="Portainer URL override")
    parser.add_argument("--username", default=None, help="Portainer username")
    parser.add_argument("--password", default=None, help="Portainer password")
    parser.add_argument("--service-id", default=None, help="Portainer service ID override")
    parser.add_argument("--endpoint-id", type=int, default=None, help="Portainer endpoint ID override")
    
    args = parser.parse_args()
    
    # Resolve parameters from arguments, env, or defaults
    url = args.url or get_env_or_default("PORTAINER_URL", DEFAULT_PORTAINER_URL)
    username = args.username or get_env_or_default("PORTAINER_USERNAME", "CortexxRick")
    password = args.password or get_env_or_default("PORTAINER_PASSWORD", "14850d130b7308@eadB")
    service_id = args.service_id or get_env_or_default("PORTAINER_SERVICE_ID", DEFAULT_SERVICE_ID)
    endpoint_id = args.endpoint_id or int(get_env_or_default("PORTAINER_ENDPOINT_ID", str(DEFAULT_ENDPOINT_ID)))
    
    new_image = f"diegosouzapw/omniroute:{args.version}"
    
    # 1. Login
    try:
        token = authenticate(url, username, password)
    except Exception as e:
        print(f"[!] Authentication failed: {e}")
        sys.exit(1)
        
    # 2. Get current running container ID to track rolling deploy
    old_container_id = None
    try:
        containers = get_containers(url, token, endpoint_id)
        for c in containers:
            names = c.get("Names", [])
            name = names[0] if names else ""
            if "cortexx-omniroute_cortexx-omniroute" in name and c.get("State") == "running":
                old_container_id = c.get("Id")
                print(f"[*] Found current running container: {name} ({old_container_id[:12]})")
                break
    except Exception as e:
        print(f"[!] Warning: could not determine old container ID: {e}")
        
    # 3. Fetch current service config
    try:
        print(f"[*] Fetching config for service {service_id}...")
        service = get_service(url, token, service_id, endpoint_id)
    except Exception as e:
        print(f"[!] Failed to get service: {e}")
        sys.exit(1)
        
    version_index = service.get("Version", {}).get("Index")
    spec = service.get("Spec")
    
    # 4. Modify Spec for Update
    task_template = spec.setdefault("TaskTemplate", {})
    task_template["ForceUpdate"] = task_template.get("ForceUpdate", 0) + 1
    
    container_spec = task_template.setdefault("ContainerSpec", {})
    old_image = container_spec.get("Image")
    print(f"[*] Current image in spec: {old_image}")
    
    container_spec["Image"] = new_image
    print(f"[*] Setting target image in spec: {new_image}")
    
    # 5. Trigger update
    try:
        print(f"[*] Sending update command to Portainer...")
        result = update_service(url, token, service_id, endpoint_id, version_index, spec)
        print(f"[+] Update response: {result}")
    except Exception as e:
        print(f"[!] Failed to update service: {e}")
        sys.exit(1)
        
    # 6. Monitor deployment
    success = monitor_deployment(url, token, endpoint_id, old_container_id)
    if not success:
        sys.exit(1)

if __name__ == "__main__":
    main()
