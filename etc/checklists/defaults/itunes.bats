load test_helper

@test "Device backups encrypted by iTunes" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  shopt -s nullglob
  for plist in /Users/*/Library/Application Support/MobileSync/Backup/*/Manifest.plist; do
    [ "$(defaults read "$plist" isEncrypted; true)" = "1" ] || [ "$(defaults read "$plist" IsEncrypted; true)" = "1" ]
  done
}
