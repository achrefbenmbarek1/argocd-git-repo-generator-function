# argo-cd-git-repo-generator-function

A composition function that generates an Argo CD Helm release based on repository URLs provided in a composite resource.

Each generated release:

- **Extracts Repository URLs:**  
  Reads the composite resource field `spec.repos-urls` to retrieve an array of Git repository URLs.

- **Transforms Repository Data:**  
  For each valid repository URL, the function:
  - Extracts the repository name and username.
  - Converts the repository name into a Kubernetes‑friendly (lowercase kebab-case) string (it is meant to be used with ESO so that is the reason behind the dynamic nature of assigning the secret through "set").
  - Generates set items that configure credential templates and repository entries.

- **Generates Helm Values:**  
  Constructs a JSON values object that includes:
  - Server configuration (e.g. service type, ingress settings, extra arguments).
  - Configurations for credentials and repositories based on the input data.
  
- **Creates a Helm Release Manifest:**  
  Uses the transformed values and set items to generate a Helm release resource (using provider-helm types) that:
  - References the Argo CD chart (name, repository, version).
  - Sets the target namespace and applies the generated values.
  
## Seeing it in action for a better understanding
1. **Open a new terminal and clone this repo using this command:**
```bash
git clone git@github.com:achrefbenmbarek1/argocd-git-repo-generator-function.git
```
2. **Head to the root directory of the project:**
```bash
cd argocd-git-repo-generator-function
```
3. **If you don't have go run this command (but before running this,make sure that you've followed the two first steps of the How to Use this section):**
```bash
nix shell nixpkgs#go
```
4. **Run this for testing purposes:**
```bash
 go run . --insecure --debug
```
5. **Open an other terminal and go straight to the example directory using this command(replace "pathToTheClonedProject" with the actual path to the cloned Project):**
```bash
 cd pathToTheClonedProject/example
```
6. **Run the first 3 steps of the How to use this section:**
(Refer to the How to Use section below for those steps.)
7. **If you don't have docker running and you went through step 6 just run this to start docker in the background:
```bash
sudo $(command -v dockerd) &
```
8. **Run this command to see the output and feel free to play with the xr.yaml to test the output of different inputs and understand more what is happening and compare the output with the composition.yam in the example directory:**
```bash
crossplane render xr.yaml composition.yaml functions.yaml
```
## How to use this:
Make sure that you've a Kubernetes cluster where you've deployed Crossplane and the Helm provider of Crossplane. You will also need to define the required resources and configurations to enable proper integration. These are what you need to define:
 If you don't have a kubernetes cluster and you're new to this (and you use linux or windows with wsl) just follow this(even if you're not I would still recommend this to install dependencies easily):
   1. **Install nix through running this command:**
     ```bash 
     sh <(curl -L https://nixos.org/nix/install) --daemon
     ```
   2. **Enable flakes though running this:**
     ```bash
     echo "experimental-features = nix-command flakes" > ~/.config/nix/nix.conf
     ```
   3. Create the environment with the necessary and preferable packages installed (from now on all the packages would be available in the shell session that you ran this command on, the packages will be unaccessible in any other session unless you run this command on it, one more thing the first time that you're gonna run this it will take a little bit of time but your subsequent runs will be basically instantanious cause the packages will be cached, at the end of this I will show you how to remove that cache and uninstall everything that we've installed with this command with another one command):
     ```bash
     nix develop ./nix-flake/kuberFlake#kuber --impure
     ```
   4. Deploy crossplane to the cluster through the run of this command from the root directory(I'm assuming that you already have a kubernetes cluster if not just perform step 7 of the "Seeing it in action for a better understanding" section and then spawn a cluster with kind which I already installed it for you in that case just run: kind create cluster --name kuber-test):
     ```bash
     helmfile sync
     ```

   5. Deploy the Crossplane Helm provider by using the provided package from the Crossplane registry and the other necessary resources.
   You can use the provided example just go under the root directory and run this 2 times: 
   ```bash
   kubectl apply -f helm-provider-stuff
   ```

   6. Deploy the function to your cluster through the use of the file in this repo
   ```bash
   kubectl apply -f argocd/argocd-git-repo-generator-function.yaml
   ```
   7. Then incorporate it in your composition, you can look the directory: "argocd", for an example of such composition (in that example I show what inputs are needed on the composition and an example claim and an xrd)
