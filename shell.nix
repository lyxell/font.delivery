{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/nixos-unstable.tar.gz") {} }:

pkgs.mkShell {
  nativeBuildInputs = [
	pkgs.go
	pkgs.gofumpt
	pkgs.gopls
	pkgs.harfbuzz
	pkgs.just
	pkgs.miniserve
	pkgs.nodePackages.typescript-language-server
	pkgs.nodejs_22
	pkgs.woff2
	pkgs.redocly
	pkgs.oapi-codegen
  ];
}
