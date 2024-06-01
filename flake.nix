{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };
  outputs = { self, nixpkgs, ... }@inputs: let
    pkgs = nixpkgs.legacyPackages.x86_64-linux;
  in rec {
    devShells.x86_64-linux.default = pkgs.mkShell {
      packages = with pkgs; [
        go
      ];
    };  
  };
}