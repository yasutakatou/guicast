from ipykernel.kernelbase import Kernel
import sys
from winpty import PtyProcess
import re

ansi_escape = re.compile(r'\x1B\[[0-?]*[ -/]*[@-~]')
proc = PtyProcess.spawn('guicast')

class EchoKernel(Kernel):
    implementation = 'Echo'
    implementation_version = '1.0'
    language = 'no-op'
    language_version = '0.1'
    language_info = {
        'name': 'Any text',
        'mimetype': 'text/plain',
        'file_extension': '.txt',
    }
    banner = "Echo kernel - as useful as a parrot"

    def do_execute(self, code, silent, store_history=True, user_expressions=None,
                   allow_stdin=False):
        if not silent:
            proc.write(code + '\n')
            rflag = True
            buffer = ''
            while rflag:
                  str = proc.readline()
                  if str.find(' - - - - - -') != -1:
                     rflag = False
                  else:
                     buffer = buffer + ansi_escape.sub('', str)
            stream_content = {'name': 'stdout', 'text': buffer }
            self.send_response(self.iopub_socket, 'stream', stream_content)

        return {'status': 'ok',
                # The base class increments the execution count
                'execution_count': self.execution_count,
                'payload': [],
                'user_expressions': {},
               }
