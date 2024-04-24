{ buildGoModule }:

buildGoModule rec {
  pname = "games";
  src = ./.;
  version = "0.1";
}