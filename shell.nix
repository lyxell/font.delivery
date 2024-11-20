{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/nixos-unstable.tar.gz") {} }:

pkgs.mkShell {
  nativeBuildInputs = [
	# Building
	pkgs.go
	pkgs.gopls
	pkgs.harfbuzz
	pkgs.just
	pkgs.tailwindcss
	pkgs.woff2
	# Development
	pkgs.gofumpt
	pkgs.miniserve
	pkgs.nodePackages.prettier
	pkgs.typescript
	pkgs.typescript-language-server
  ];
}
