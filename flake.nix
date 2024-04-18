{
  description = "Pacur - Automated deb, rpm and pkgbuild build system";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    nixpkgs,
    flake-utils,
    self,
  } @ inputs:
    with inputs;
      flake-utils.lib.eachDefaultSystem (
        system: let
          pkgs = nixpkgs.legacyPackages.${system};
        in rec {
          packages.default = pkgs.callPackage ./default.nix {};
        }
      );
}
