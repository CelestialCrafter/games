{
  description = "games";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };
  outputs = { self, nixpkgs, ... }@inputs: let
    pkgs = nixpkgs.legacyPackages.x86_64-linux;
  in rec {
    packages.x86_64-linux.games = pkgs.callPackage ./default.nix {};
    packages.x86_64-linux.default = packages.x86_64-linux.games;

    devShells.x86_64-linux.default = pkgs.mkShell {
      packages = with pkgs; [
        go
      ];
    };  
  };
}