
allow_k8s_contexts(os.getenv("TILT_ALLOW_CONTEXT"))

k8s_yaml('example-http.yaml')
k8s_yaml('example-exec.yaml')

k8s_resource('http-check', port_forwards=['127.0.0.1:9090:80', '127.0.0.1:8080:8080'])
k8s_resource('exec-check', port_forwards=['127.0.0.1:9091:80'])

target='prod'
live_update=[]
if os.environ.get('PROD', '') ==  '':
  target='build-env'
  live_update=[
    sync('pkg',    '/app/pkg'),
    sync('cmd',    '/app/cmd'),
    sync('go.mod', '/app/go.mod'),
    sync('go.sum', '/app/go.sum'),
    run('go install ./cmd/network-health-client'),
    run('go install ./cmd/network-health-server'),
  ]

docker_build(
  ref='network-health-image',
  context='.',
  live_update=live_update,
  target=target,
  only=[ 'go.mod'
       , 'go.sum'
       , 'pkg'
       , 'cmd'
       , 'entrypoint.sh'
  ],
  ignore=[ '.git'
         , '*/*_test.go'
         , '*.yaml'
  ],
)
