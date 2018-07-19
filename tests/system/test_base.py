from akkibeat import BaseTest

import os


class Test(BaseTest):

    def test_base(self):
        """
        Basic test with exiting Akkibeat normally
        """
        self.render_config_template(
            path=os.path.abspath(self.working_dir) + "/log/*"
        )

        akkibeat_proc = self.start_beat()
        self.wait_until(lambda: self.log_contains("akkibeat is running"))
        exit_code = akkibeat_proc.kill_and_wait()
        assert exit_code == 0
