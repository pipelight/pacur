{
  pkgs ? import <nixpkgs> {},
  lib,
  ...
}: pkgs.buildGoModule rec {
  pname = "pacur";
  version = "master";
    src = ./.;
  }
