
allow_k8s_contexts(os.getenv("TILT_ALLOW_CONTEXT"))

k8s_yaml('example-http.yaml')
k8s_yaml('example-exec.yaml')

k8s_resource('http-check', port_forwards=['127.0.0.1:9090:80', '127.0.0.1:8080:8080'])
k8s_resource('exec-check', port_forwards=['127.0.0.1:9091:80'])

docker_build(
  ref='network-health-image',
  context='.',
  ignore=[
    'vendor',
    'kubernetes.yaml',
  ],
)
