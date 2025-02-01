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
      ko
];
NIX_DEV="kuber-env";

    shellHook = ''
            if [ -f ~/.zshrc ]; then
              zsh
            elif [ -f ~/.bashrc ]; then
              bash
            fi

          '';
    };
  };
}


