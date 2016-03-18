
import os
import os.path
import shutil
import sys
import json
import tarfile

import requests

from third_party.py import gflags

gflags.DEFINE_string(
    'registry', None,
    'foo')
gflags.MarkFlagAsRequired('registry')

gflags.DEFINE_string(
    'repository', None,
    'foo')
gflags.MarkFlagAsRequired('repository')

gflags.DEFINE_string(
    'image', None,
    'foo')
gflags.MarkFlagAsRequired('image')

gflags.DEFINE_string(
    'digest', None,
    'SHA256, SHA384, SHA512')

gflags.DEFINE_string(
    'out_path', None,
    'SHA256, SHA384, SHA512')
gflags.MarkFlagAsRequired('digest')

FLAGS = gflags.FLAGS

verify='/etc/ssl/certs/ca-certificates.crt'

def main(unused_argv):
  requests.packages.urllib3.disable_warnings()
  r = requests.get("https://%s/v2/%s/%s/manifests/%s" % (FLAGS.registry, FLAGS.repository, FLAGS.image, FLAGS.digest), verify=verify)
  manifest = r.json()
  shutil.rmtree("tmp/", ignore_errors=True)
  os.mkdir("tmp/")
  for i in range(0, len(manifest["fsLayers"])):
    layer_json = manifest["history"][i]["v1Compatibility"]
    layer = json.loads(layer_json)
    layer_path = "tmp/%s/" % layer["id"]
    os.mkdir(layer_path)
    with open(layer_path + "json", 'w') as layer_json_file:
      layer_json_file.write(layer_json)
    with open(layer_path + "VERSION", 'w') as layer_version_file:
      layer_version_file.write("1.0")
    with open(layer_path + "layer.tar", 'w') as layer_tar_file:
      r = requests.get("https://%s/v2/%s/%s/blobs/%s" % (FLAGS.registry, FLAGS.repository, FLAGS.image, manifest["fsLayers"][i]["blobSum"]), verify=verify)
      for chunk in r.iter_content(chunk_size=1024):
        if chunk:
          layer_tar_file.write(chunk)
  tar = tarfile.open(FLAGS.out_path, mode="w|")
  for dir_name, subdir_list, file_list in os.walk("tmp/"):
    for filename in file_list:
      path = os.path.join(dir_name, filename)
      info = tarfile.TarInfo(path[len("tmp/"):])
      info.size = os.stat(path).st_size
      with open(path, 'r') as f:
        tar.addfile(info, f)
  tar.close()

if __name__ == '__main__':
  main(FLAGS(sys.argv))
