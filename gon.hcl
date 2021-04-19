source = ["./dist/datree-macos_darwin_amd64/datree"]
bundle_id = "io.datree"

apple_id {
  username = "yishay@datree.io"
  password = "@env:AC_PASSWORD"
}

sign {
  application_identity = "Developer ID Application: Datree Group Inc"
}

zip {
  output_path = "./dist/datree-macos.zip"
}
