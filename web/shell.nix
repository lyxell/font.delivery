{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/nixos-unstable.tar.gz") {} }:

pkgs.mkShell {
  nativeBuildInputs = [
	pkgs.nodePackages.typescript-language-server
	pkgs.nodejs_22
	pkgs.just
  ];
}
