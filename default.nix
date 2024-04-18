{
  pkgs ? import <nixpkgs> {},
  lib,
  ...
}:
pkgs.buildGoModule rec {
  pname = "pacur";
  version = "master";
  src = ./.;
  # vendorHash = lib.fakeHash;
  vendorHash = "sha256-jJQYwQxcXRrmK4idnynGhIjQzGOp44Lj3EEuXMPrPII=";
  # buildInputs = with pkgs; [podman];
}
