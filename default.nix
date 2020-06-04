{ stdenv, pkgs }:

stdenv.mkDerivation {
  name = "gomodoro";
  src = ./.;

  buildInputs = [ pkgs.go ];

  buildPhase = "GOCACHE=/tmp go build gomodoro.go";

  installPhase = ''
    mkdir -p $out/bin;
    cp gomodoro $out/bin/
   '';
 }
