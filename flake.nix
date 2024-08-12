{
  description = "Game engine dev environemnt";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }: {
    devShell.x86_64-linux = nixpkgs.legacyPackages.x86_64-linux.mkShell {
      buildInputs = with nixpkgs.legacyPackages.x86_64-linux; [
        go
        libGL
        libGLU
        glfw
        xorg.libX11
        xorg.libXi
        xorg.libXcursor
        xorg.libXrandr
        xorg.libXinerama
        xorg.libXxf86vm
        pkg-config
        assimp
      ];
    };
  };
}

