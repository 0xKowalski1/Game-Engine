{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.cope
    pkgs.go
    pkgs.libGL
    pkgs.libGLU
    pkgs.glfw
    pkgs.xorg.libX11
    pkgs.xorg.libXi
    pkgs.xorg.libXcursor
    pkgs.xorg.libXrandr
    pkgs.xorg.libXinerama
    pkgs.xorg.libXxf86vm
    pkgs.pkg-config  
  ];


}

