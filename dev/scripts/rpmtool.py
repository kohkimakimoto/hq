#!/usr/bin/env python
from __future__ import division, print_function, absolute_import, unicode_literals
import argparse, os, sys, re, subprocess, textwrap, shutil, tempfile, grp

# utilities for compatibility.
PY2 = sys.version_info[0] == 2
PY3 = sys.version_info[0] == 3

if PY2:
    input = raw_input
    def as_bytes(s, encoding='utf-8'):
        if isinstance(s, str):
            return s
        else:
            return s.encode(encoding)

    def as_string(s, encoding='utf-8'):
        if isinstance(s, unicode):
            return s
        else:
            return s.decode(encoding)
else:
    input = input
    def as_bytes(s, encoding='utf8'):
        if isinstance(s, bytes):
            return s
        else:
            return s.encode(encoding)

    def as_string(s, encoding='utf8'):
        if isinstance(s, str):
            return s
        else:
            return s.decode(encoding)

def shell_escape(s):
    return "'" + s.replace("'", "'\"'\"'") + "'"

class TemporaryDirectory(object):
    def __init__(self, suffix="", prefix="tmp", dir=None):
        self.name = tempfile.mkdtemp(suffix, prefix, dir)

    def __enter__(self):
        return self.name

    def __exit__(self, exc, value, tb):
        self.cleanup()

    def cleanup(self):
        shutil.rmtree(self.name)

def platform_to_image(platform):
    # currently, we only support el7
    if platform == "el7":
        return "centos:7"
    else:
        raise ValueError("unknown platform: %s" % (platform))

def build_image(platform, forceBuild=False):
    base_image = platform_to_image(platform)
    image = "kohkimakimoto/rpmtool:%s" % platform
    image_id = subprocess.check_output("docker images -q %s" % (image),  shell=True).strip()
    if forceBuild or image_id == "":
        # Create temporary directory to store Dockerfile.
        with TemporaryDirectory(prefix="rpmtool-") as datadir:
            with open(os.path.join(datadir, "Dockerfile"), 'w') as f:
                f.write(textwrap.dedent('''\
                    FROM %s
                    RUN yum -y install epel-release \\
                    && yum clean all \\
                    && yum -y --setopt=tsflags='' install \\
                        gcc make rpmdevtools mock perl sudo tar zlib-devel createrepo

                    RUN adduser --comment "RPM Tool Bot" --home /home/rpmtoolbot --create-home rpmtoolbot
                    RUN sed -i -e 's/Defaults    env_reset/#Defaults    env_reset/' /etc/sudoers \
                        && sed -i -e 's/Defaults    secure_path/#Defaults    secure_path/' /etc/sudoers \
                        && echo 'Defaults    env_keep += "PATH"' >> /etc/sudoers \
                        && echo 'rpmtoolbot    ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

                    USER rpmtoolbot
                    RUN rpmdev-setuptree
                    WORKDIR /home/rpmtoolbot
                    CMD ["/bin/bash"]
                ''' % base_image))
            # build image
            subprocess.check_call("cd %s && docker build -t %s ." % (datadir, image),  shell=True)
    return image

def build_command(args):
    image = build_image(args.platform, args.build if args.build else False)
    spec_file = args.spec_file
    rpmbuild_dir = os.path.abspath(os.path.dirname(os.path.dirname(spec_file)))
    out_dir = args.out_dir

    pre_commands = args.pre
    post_commands = args.post

    try:
        # run pre command
        for cmd in pre_commands:
            subprocess.check_call(cmd, shell=True)

        entryscript = textwrap.dedent(shell_escape("""\
            set -e
            cp -pr /shared/rpmbuild/SPECS $HOME/rpmbuild/
            cp -pr /shared/rpmbuild/SOURCES $HOME/rpmbuild/

            cd $HOME
            rpmbuild -ba ~/rpmbuild/SPECS/%(spec_file)s
            cp -pr $HOME/rpmbuild/RPMS/* /shared/outputs/
            cp -pr $HOME/rpmbuild/SRPMS/* /shared/outputs/
        """ % {"spec_file": os.path.basename(spec_file)}))
        subprocess.check_call("docker run --rm -v %s:/shared/outputs -v %s:/shared/rpmbuild %s /bin/bash -c %s" % (out_dir, rpmbuild_dir, image, entryscript),  shell=True)
    finally:
        # run post command
        for cmd in post_commands:
            subprocess.check_call(cmd, shell=True)

def createrepo_command(args):
    image = build_image(args.platform, args.build if args.build else False)
    createrepo_args = args.args
    subprocess.check_call("docker run --rm -u \"$(id -u):$(id -g)\" -v $PWD:/build -w /build %s createrepo %s" % (image, " ".join(createrepo_args)),  shell=True)


def main():
    parser = argparse.ArgumentParser(
        description="A portable tool for building rpm package by using docker container",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=textwrap.dedent('''\
            examples:
                rpmtool.py build path/to/xxx.spec

            Copyright (c) Kohki Makimoto <kohki.makimoto@gmail.com>
            The MIT License (MIT)
        '''))

    subparsers = parser.add_subparsers(title='subcommands')

    parser_build = subparsers.add_parser('build', help='build a RPM package')
    parser_build.add_argument("-b", "--build", dest="build", action="store_true", help="force to build an image before starting the container.")
    parser_build.add_argument("--platform", dest="platform", metavar="PLATFORM", help="specify the platform. (default: 'el7')", default="el7")
    parser_build.add_argument("--out", dest="out_dir", metavar="DIR", help="specify the output directory", default=os.getcwd())
    parser_build.add_argument("--pre", dest="pre", metavar="COMMAND", nargs='*', help="commands executed before building", default=[])
    parser_build.add_argument("--post", dest="post", metavar="COMMAND", nargs='*',help="commands executed after building", default=[])
    parser_build.add_argument('spec_file', help="the RPM spec file path.")
    parser_build.set_defaults(func=build_command)

    parser_createrepo = subparsers.add_parser('createrepo', help='run createrepo command in a docker container')
    parser_createrepo.add_argument("-b", "--build", dest="build", action="store_true", help="force to build an image before starting the container.")
    parser_createrepo.add_argument("--platform", dest="platform", metavar="PLATFORM", help="specify the platform. (default: 'el7')", default="el7")
    parser_createrepo.add_argument("args", metavar="ARGS", nargs='*', help="arguments of createrepo command", default=[])
    parser_createrepo.set_defaults(func=createrepo_command)

    if len(sys.argv) == 1:
        parser.print_help()
        sys.exit(1)

    args = parser.parse_args()
    args.func(args)

if __name__ == '__main__': main()
