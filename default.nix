{
  lib,
  pkgs,
  ...
}:
pkgs.buildGoModule (finalAttrs: {
  pname = "anilistarr";
  version = "0.3.0";
  src =
    with pkgs.lib.fileset;
    toSource {
      root = ./.;
      fileset = intersection (gitTracked ./.) (unions [
        (fileFilter (file: file.hasExt "go") ./.)
        ./cmd/handler
        ./go.mod
        ./go.sum
        ./internal
        ./pkg
        ./sqlc.yaml
        ./swagger.yaml
      ]);
    };
  vendorHash = "sha256-SEYA4jNIUF0qV/gtQ2w80R0D9iUDoIPSqIysIJ9PZRs=";
  modules = ./gomod2nix.toml;
  subPackages = [ "cmd/handler" ];
  ldflags = [
    "-s"
    "-w"
    "-X main.version=${finalAttrs.version}"
  ];
  meta = {
    description = "anilist custom list provider for sonarr/radarr";
    homepage = "https://github.com/wwmoraes/anilistarr";
    license = lib.licenses.mit;
    maintainers = [ lib.maintainers.wwmoraes ];
    mainProgram = "handler";
  };
})
