load test_helper

@test "File sharing disabled" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  run defaults read /Library/Preferences/SystemConfiguration/com.apple.smb.server.plist EnabledServices
  [ $status -ne 0 ] || for w in $output; do [[ "$w" = "disk" ]]; done
}
