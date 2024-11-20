{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/nixos-unstable.tar.gz") {} }:

pkgs.mkShell {
  nativeBuildInputs = [
	pkgs.go
	pkgs.gofumpt
	pkgs.gopls
	pkgs.harfbuzz
	pkgs.just
	pkgs.miniserve
	pkgs.typescript
	pkgs.typescript-language-server
	pkgs.woff2
  ];
}
