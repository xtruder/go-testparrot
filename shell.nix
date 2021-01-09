{ src ? builtins.fetchTarball "https://github.com/NixOS/nixpkgs/archive/nixos-20.09.tar.gz",
  pkgs ? import src {}}:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go_1_15
    gopls
    delve
    go-outline
  ];

  hardeningDisable = [ "all" ];

  GO111MODULE = "on";
}
