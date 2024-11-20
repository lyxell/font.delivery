{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/nixos-unstable.tar.gz") {} }:

pkgs.mkShell {
  nativeBuildInputs = [
	pkgs.go
	pkgs.gofumpt
	pkgs.gopls
	pkgs.harfbuzz
	pkgs.just
	pkgs.protobuf
	pkgs.protoc-gen-go
	pkgs.woff2
  ];
}
