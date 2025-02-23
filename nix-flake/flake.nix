{
  description = "kubernetes environment";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };
  outputs = { self, nixpkgs }: 
  let
  system = "x86_64-linux";
  pkgs = nixpkgs.legacyPackages.x86_64-linux;
  in
  {
    devShells.x86_64-linux.kuber = pkgs.mkShell  {
      buildInputs = with pkgs; [ 
      kind
      docker
      kubectl
      helmfile
      (pkgs.wrapHelm pkgs.kubernetes-helm { plugins = [ pkgs.kubernetes-helmPlugins.helm-diff ]; })
      helm-ls
      argocd
      crossplane-cli
];
NIX_DEV="kuber-env";

    shellHook = ''
            if [ -f ~/.zshrc ]; then
              zsh
            elif [ -f ~/.bashrc ]; then
              bash
            fi

            # Start Docker daemon in the background if it's not already running
            # if ! pgrep -x dockerd > /dev/null; then
            #   echo "Starting Docker daemon..."
            #   sudo $(command -v dockerd)
            # fi
          '';
    };
  };
}


