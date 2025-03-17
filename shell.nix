let
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/archive/5d9b5431f967007b3952c057fc92af49a4c5f3b2.tar.gz";

  pkgs = import nixpkgs {
    config = { };
    overlays = [ ];
  };

in

pkgs.mkShellNoCC {
  packages = with pkgs; [
    go
    tailwindcss_4
    templ
  ];

  TAILWIND_CLI = "tailwindcss";
}
