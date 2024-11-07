import sys
import subprocess
from pathlib import Path


def generate_proto_files():
    try:
        import grpc_tools.protoc  # noqa
    except ImportError:
        print("Installing required dependencies...")
        subprocess.check_call([sys.executable, "-m", "pip", "install", "grpcio-tools"])

    # Get the directory containing the .proto files
    proto_dir = Path(__file__).parent.absolute()

    # Proto files to generate
    proto_files = [
        "common.proto",
        "vector.proto",
        "graph.proto",
    ]  # Order matters! Dependencies first

    for proto_file in proto_files:
        proto_path = proto_dir / proto_file
        if not proto_path.exists():
            print(f"Error: Proto file not found: {proto_path}")
            continue

        print(f"Generating Python files for {proto_file}...")

        # Command to generate Python code from proto file
        cmd = [
            sys.executable,
            "-m",
            "grpc_tools.protoc",
            f"--proto_path={proto_dir}",
            f"--python_out={proto_dir}",
            f"--grpc_python_out={proto_dir}",
            "--experimental_allow_proto3_optional",  # Add this flag
            str(proto_path),
        ]

        try:
            # Run the protoc compiler
            result = subprocess.run(cmd, capture_output=True, text=True)
            if result.returncode != 0:
                print(f"Error output: {result.stderr}")
                raise subprocess.CalledProcessError(
                    result.returncode, cmd, result.stdout, result.stderr
                )
            print(f"Successfully generated files for {proto_file}")

            # Fix imports in generated files
            generated_file = proto_dir / f"{proto_file.replace('.proto', '_pb2.py')}"
            if generated_file.exists():
                content = generated_file.read_text()
                content = content.replace(
                    "import common_pb2", "from . import common_pb2"
                )
                content = content.replace(
                    "import vector_pb2", "from . import vector_pb2"
                )
                generated_file.write_text(content)

            grpc_file = proto_dir / f"{proto_file.replace('.proto', '_pb2_grpc.py')}"
            if grpc_file.exists():
                content = grpc_file.read_text()
                content = content.replace(
                    "import common_pb2", "from . import common_pb2"
                )
                content = content.replace(
                    "import vector_pb2", "from . import vector_pb2"
                )
                content = content.replace("import graph_pb2", "from . import graph_pb2")
                grpc_file.write_text(content)

        except subprocess.CalledProcessError as e:
            print(f"Error generating files for {proto_file}:")
            print(f"Command output: {e.output}")
            print(f"Error output: {e.stderr}")
            raise


if __name__ == "__main__":
    try:
        generate_proto_files()
    except Exception as e:
        print(f"Error: {str(e)}")
        sys.exit(1)
