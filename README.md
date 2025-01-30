# argo-cd-git-repo-generator-function

A [composition function][functions] in [Go][go] that names the repositories, the credentialTemplates, and the name keys in the Helm release resource of ArgoCD provided by the Crossplane Helm provider, dynamically based on an array of repo URLs given to the Crossplane claim.

## How to use this:

Make sure that you've a Kubernetes cluster where you've deployed Crossplane and the Helm provider of Crossplane. You will also need to define the required resources and configurations to enable proper integration. These are what you need to define:

1. **Install Crossplane and the Helm Provider**:
   - Deploy Crossplane in your Kubernetes cluster if it’s not already installed.
   If you don't have a kubernetes cluster and you're new to this (and you use linux or windows with wsl) just follow this:
     - install nix through running this command:
     ```bash 
     sh <(curl -L https://nixos.org/nix/install) --daemon
     ```
     - enable flakes though running this:
     ```bash
     echo "experimental-features = nix-command flakes" > ~/.config/nix/nix.conf
     ```
     - Create the environment with the necessary and preferable packages installed (from now on all the packages would be available in the shell session that you ran this command on, the packages will be unaccessible in any other session unless you run this command on it, one more thing the first time that you're gonna run this it will take a little bit of time but your subsequent runs will be basically instantanious cause the packages will be cached, at the end of this I will show you how to remove that cache and uninstall everything that we've installed with this command with another one command):
     ```bash
     nix develop ./nix-flake/kuberFlake#kuber --impure
     ```
   - Deploy crossplane to the cluster through the run of this command from the root directory:
     ```bash
     helmfile sync
     ```

   - Deploy the Crossplane Helm provider by using the provided package from the Crossplane registry and the other necessary resources.
   You can use the provided example just go under the root directory and run this 2 times: 
   ```bash
   kubectl apply -f helm-provider-stuff
   ```

