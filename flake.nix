{
  description = "Monorepo";

  inputs = {
    # nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    nixpkgs.url = "github:NixOS/nixpkgs/22779946778dabdc33e294970cd80b6e41aa8192";
    systems.url = "github:nix-systems/default";
    devenv.url = "github:cachix/devenv/34e6461fd76b5f51ad5f8214f5cf22c4cd7a196e";
  };

  nixConfig = {
    extra-trusted-public-keys = "devenv.cachix.org-1:w1cLUi8dv3hnoSPGAuibQv+f9TZLr6cv/Hm9XgU50cw=";
    extra-substituters = "https://devenv.cachix.org";
  };

  outputs = { self, nixpkgs, devenv, systems, ... } @ inputs:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forEachSystem = nixpkgs.lib.genAttrs (import systems);
    in
    {
      devShells = forEachSystem
        (system:
          let
            pkgs = nixpkgs.legacyPackages.${system};
          in
          {
            default = devenv.lib.mkShell {
              inherit inputs pkgs;
              modules = [
                {
                  # https://devenv.sh/reference/options/
                  # services = {};

                  languages.rust.enable = true;
                  languages.go.enable = true;

                  packages = with pkgs; [
                    git
                    workshop-runner
                    wrangler
                  ];

                  enterShell = ''
                  '';
                }
              ];
            };
          });
    };
}
