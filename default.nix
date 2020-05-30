with import <nixpkgs> {};

stdenv.mkDerivation {
  name = "gomodoro";
  src = ./.;

  buildInputs = [ go ];
}
