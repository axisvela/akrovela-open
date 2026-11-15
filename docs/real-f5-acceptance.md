# F5 Real-Run Acceptance Procedures

Status: external verification procedure. This document describes how to run the
real ROS 2/Nav2, Gazebo and PEIRAVELA paths. It is not executed evidence, not a
production deployment, and not a safety or physical-effect claim. `akrovela-open`
publishes the contract and this procedure; the core implementation lives in the
`akrovela` repository.

## Common Setup

Recommended host: Ubuntu 24.04 x86_64 with ROS 2 Jazzy, Gazebo Harmonic, Nav2 and
TurtleBot 4 simulation (or an equivalent Nav2 robot).

Required:

- Go 1.25.5 or newer.
- Python 3.10 or newer with `rclpy`, `nav2_msgs`, `geometry_msgs` and the Nav2
  action server.
- PostgreSQL 17 and `psql`.
- A disposable database and disposable tokens. Never commit credentials.
- Loopback-only bridge and verifier listeners.

Build and migrate the core:

```bash
git clone https://github.com/axisvela/akrovela
cd akrovela/backend
go build ./...
export AKROVELA_TEST_DATABASE_URL='postgres://user:password@127.0.0.1:5432/akrovela_accept?sslmode=disable'
go run ./cmd/akrovela-dbadmin up
```

The migration set is `0001`-`0014`. The application never applies migrations at
startup. Record the commit, Go version, ROS distribution and Gazebo version.

## Track A: Real ROS 2/Nav2 + rclpy

### A1. Start the simulator and Nav2

```bash
source /opt/ros/jazzy/setup.bash
# Start the chosen TurtleBot 4 / Nav2 simulation, for example:
ros2 launch turtlebot4_gz_bringup turtlebot4_gz.launch.py
```

In a second shell:

```bash
source /opt/ros/jazzy/setup.bash
ros2 action list | grep navigate_to_pose
ros2 action info /navigate_to_pose
ros2 topic echo /odom --once
```

Record the action name, map frame, robot frame, Gazebo world and initial pose.

### A2. Start the AKROVELA ROS 2 bridge with the real backend

From `akrovela/backend`:

```bash
export AKROVELA_ROS2_BRIDGE_BIND=127.0.0.1
export AKROVELA_ROS2_BRIDGE_PORT=2055
export AKROVELA_ROS2_BRIDGE_TOKEN='replace-with-disposable-token'
export AKROVELA_ROS2_BRIDGE_BACKEND=rclpy
export AKROVELA_ROS2_BRIDGE_STATE_FILE=/var/lib/akrovela/ros2-bridge-state.json
python -m ros2bridge
```

`rclpy` must connect to the Nav2 `/navigate_to_pose` action server. The bridge
returns execution facts only; it never returns arrival, `effect_status`,
`stop_confirmed` or Mission state.

### A3. Smoke-test the bridge

Use the goal values from the actual simulation map:

```bash
export BRIDGE_TOKEN='replace-with-disposable-token'
curl -sS -X POST http://127.0.0.1:2055/v1/submit \
  -H "Authorization: Bearer $BRIDGE_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"operation_key":"op:real-1","attempt_id":"attempt:real-1","execution_epoch":7,"fencing_token":42,"valid_until":"2026-11-15T12:05:00Z","goal":{"x":12.4,"y":-3.1,"yaw":1.5708,"frame_id":"map"},"position_tolerance_m":0.25,"yaw_tolerance_rad":0.15,"behavior_tree":"navigate_to_pose_w_replanning_and_recovery.xml","planner_id":"GridBased","controller_id":"FollowPath","max_duration_ms":5000}'

curl -sS "http://127.0.0.1:2055/v1/status?operation_key=op:real-1" \
  -H "Authorization: Bearer $BRIDGE_TOKEN"

curl -sS -X POST http://127.0.0.1:2055/v1/cancel \
  -H "Authorization: Bearer $BRIDGE_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"operation_key":"op:real-1","cancel_request_id":"cancel:real-1","deadline_ms":500}'
```

Record the accepted/terminal states, the goal id and any cancellation fact. A
cancel response is not a stop confirmation.

### A4. Start the arrival verifier

The live Gazebo source is not implemented in the current release. For this ROS
2/Nav2 track use either:

- `AKROVELA_ARRIVAL_SIMULATED_SOURCE=true` with a pose matching the submitted
  goal (`AKROVELA_ARRIVAL_SIMULATED_X`, `_Y`, `_YAW`), or
- `AKROVELA_ARRIVAL_OBSERVATION_FILE` with recorded observations and
  `AKROVELA_ARRIVAL_TRANSFORM_FILE` with `"validated": true`.

```bash
export AKROVELA_ARRIVAL_VERIFIER_BIND=127.0.0.1
export AKROVELA_ARRIVAL_VERIFIER_PORT=2054
export AKROVELA_ARRIVAL_VERIFIER_TOKEN='replace-with-disposable-token'
export AKROVELA_ARRIVAL_SIMULATED_SOURCE=true
export AKROVELA_ARRIVAL_SIMULATED_X=12.4
export AKROVELA_ARRIVAL_SIMULATED_Y=-3.1
export AKROVELA_ARRIVAL_SIMULATED_YAW=1.5708
./akrovela-arrival-verifier
```

The simulated source is development-only and is not Gazebo physics evidence.

### A5. Start the AKROVELA API with the Nav2 executor

Use `deploy/profiles/ros2-gazebo.env.example` as the template, replacing every
placeholder. The API client variable is `AKROVELA_VERIFIER_URL` (not
`AKROVELA_ARRIVAL_VERIFIER_URL`).

```bash
export AKROVELA_DATABASE_URL='postgres://user:password@127.0.0.1:5432/akrovela_accept?sslmode=disable'
export AKROVELA_LISTEN_ADDR=127.0.0.1:2046
export AKROVELA_API_TOKEN='replace-with-disposable-token'
export AKROVELA_TOKEN_TENANT=tenant-accept
export AKROVELA_TOKEN_WORKSPACE=wh3
export AKROVELA_TOKEN_SITE=site-wh3
export AKROVELA_NODE_TENANT=tenant-accept
export AKROVELA_NODE_WORKSPACE=wh3
export AKROVELA_NODE_SITE=site-wh3
export AKROVELA_NODE_ID=akrovela-node-wh3-01
export AKROVELA_NODE_BOOT_EPOCH=7
export AKROVELA_EXECUTOR=nav2
export AKROVELA_NAV2_BRIDGE_URL=http://127.0.0.1:2055
export AKROVELA_NAV2_BRIDGE_TOKEN="$BRIDGE_TOKEN"
export AKROVELA_VERIFIER_URL=http://127.0.0.1:2054
export AKROVELA_VERIFIER_TOKEN="$AKROVELA_ARRIVAL_VERIFIER_TOKEN"
./akrovela-api
```

Submit one operation. The recommended external path is the reviewed
KINETOVELA delegated-operation contract; the existing cross-product E2E script
can be used as an orchestration reference, replacing its fake bridge with the
real bridge above. A direct `POST /v1/operations` request must carry matching
RFC 8785 JCS payload and binding digests; use the reviewed fixtures and SDK
types rather than hand-computing them in a shell.

### A6. Acceptance checklist

- The bridge reaches a real `rclpy` goal and returns execution facts.
- The bridge rejects unknown fields, bad tolerances, expired validity and wrong
  bearer tokens.
- The API receipt records `adapter=ros2-nav2` and the bridge goal id.
- `effect_status=verified` only when the independent verifier returns
  `verified` with an accepted bound observation; otherwise it stays `unknown`.
- Cancellation is recorded as a request; no receipt claims `stop_confirmed`.
- No ROS topic/action result is treated as physical-effect proof.

Attach: bridge logs, `ros2 action info`, API query/receipt JSON, bridge state
file, verifier response and the exact goal/map/frame values. Sanitize tokens.

## Track B: Live Gazebo Observation Source

Current gap: `GazeboSource` returns `ErrSourceUnavailable`. A live source must
implement the observation contract with these fields per observation:
`ObservationID`, `SourceID`, `SourceEpoch`, `Seq`, `RobotRef`, `MapRef`,
`OperationKey`, `AttemptID`, `ExecutionEpoch`, `TargetDigest`, `ObservedAt`,
`ReceivedAt`, `X`, `Y`, `Yaw`, `Digest`.

Required behavior:

1. Subscribe to Gazebo ground-truth pose for the robot.
2. For each verifier request, return only observations matching all binding
   fields and the `[From, To]` freshness window.
3. Increment `Seq` monotonically per `SourceID` and keep `SourceEpoch` stable
   across a process restart; a reset must produce `unknown`, not a replay.
4. Timestamp authority time, not only simulator wall time; declare the
   simulation-to-authority offset if used.
5. Return a content digest for each observation.

Interim verification: use `AKROVELA_ARRIVAL_OBSERVATION_FILE` with recorded
ground truth and `AKROVELA_ARRIVAL_TRANSFORM_FILE` with `"validated": true`.
This proves the verifier path, not a live feed. The target digest and attempt
identity are operation-specific, so a static file is not a substitute for a live
source in a full run.

Acceptance: a matching fresh observation returns `verified`; stale, replayed,
wrong-source, wrong-binding and unvalidated-transform cases return
`rejected`/`unknown`; the transform is pinned and reviewed.

## Track C: PEIRAVELA robot-lab-nav2 Experiment

PEIRAVELA owns experiment and simulation acceptance. Use its quickstart and
control-plane docs.

1. Build and start PEIRAVELA:

```bash
cd peiravela/backend
go build ./cmd/api-server ./cmd/control-plane
export PEIRAVELA_ARTIFACT_ROOT="$PWD/tmp/akrovela-nav2-artifacts"
./api-server
```

2. Start the AKROVELA bridge, verifier and API from Track A. Confirm
   `AKROVELA_NAV2_BRIDGE_URL`, `AKROVELA_NAV2_BRIDGE_TOKEN` and
   `AKROVELA_VERIFIER_URL` point at the running services.

3. Confirm discovery:

```bash
curl -s http://127.0.0.1:8080/sut-adapters
```

The response must include `akrovela-nav2` with `nav2-submit`, `nav2-status` and
`nav2-cancel` capabilities.

4. Run the robot-lab Nav2 catalog:

```bash
./control-plane suite-run \
  profiles/robot-lab/suites/nav2/catalog.yaml \
  contracts/schemas/scenario-v1.schema.json \
  akrovela-nav2-navigate-to-pose-v1 benign-control \
  tmp/akrovela-nav2-artifacts
```

The catalog options carry operation/attempt identity, goal, tolerances and
Nav2 allowlists. Adjust the goal/map to the actual simulation. If the verifier
uses the simulated source, set its pose to the same goal.

5. Export and independently consume the evidence package:

```bash
./control-plane export-package tmp/akrovela-nav2-artifacts <attempt-id> tmp/akrovela-nav2-package.tar
pwsh scripts/verify-package-independent.ps1 -Package tmp/akrovela-nav2-package.tar -OutDir tmp/akrovela-nav2-consume
```

6. Acceptance checklist:
   - `GET /sut-adapters` lists `akrovela-nav2`.
   - The run records bridge submit/status/cancel and verifier observations.
   - Evidence labels preserve the AKROVELA `contract_mode` and PEIRAVELA
     `synthetic` label; a PEIRAVELA run does not complete a Mission.
   - The evidence package verifies with tar + SHA-256 only.
   - No ROS 2, Gazebo or hardware claim is made unless the external run
     actually used them.

Attach: control-plane run/attempt JSON, evidence package digest, verification
output, bridge/verifier logs and exact catalog revision.

## Reporting An External Run

File an issue in `akrovela-open` with:

- repository commit SHAs and versions;
- host OS, ROS 2, Gazebo, Nav2 and PEIRAVELA versions;
- exact commands and sanitized environment variable names;
- bridge/verifier/API/PEIRAVELA logs and exit codes;
- operation query and receipt JSON;
- evidence package digest and independent-verification output;
- explicit list of limitations and anything not run.

Do not attach tokens, credentials, tenant data, model weights or unrestricted
payloads. A pass on one host, profile or simulator does not pass another.
