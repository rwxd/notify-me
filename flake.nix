{
  description = "Small command line tool to notify myself through various services";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.notify-me = pkgs.buildGoModule {
          pname = "notify-me";
          version = self.shortRev or "dev";
          src = ./.;
          vendorHash = null;
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.notify-me}/bin/notify-me";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ go ];
        };
      }
    );
}
