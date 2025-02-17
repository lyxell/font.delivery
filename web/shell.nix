{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/4989a246d7a390a859852baddb1013f825435cee.tar.gz") {} }:

pkgs.mkShell {
  nativeBuildInputs = [
	pkgs.nodejs_22
	pkgs.just
  ];
}
