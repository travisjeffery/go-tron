load test_helper

@test "SSH has non-blank passphrase" {
  echo_url_for_description "$BATS_TEST_DESCRIPTION"
  [ -z $(grep -L ENCRYPTED "$HOME/.ssh/id_rsa") ]
}
