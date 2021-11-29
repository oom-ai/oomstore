import sys
import grpc
import logging
import time
from pathlib import Path
from subprocess import Popen

from .codegen import oomagent_pb2
from .codegen import oomagent_pb2_grpc

# Convert google.protobuf.pyext._message.MessageMapContainer object to Python dictionary
def map_container_to_dict(map_container):
    return dict({k: getattr(v, v.WhichOneof("kind")) for k, v in map_container.items()})


class Client(object):
    def __init__(self, port: int, config_path: str):
        try:
            self.oomagent = Popen(
                ["oomagent", "--config", config_path, "-p", str(port)]
            )
        except Exception as e:
            logging.error(e)
            sys.exit(1)

        self.addr = "127.0.0.1:%d" % port

        # wait for oomagent to start
        time.sleep(2)

    def __del__(self):
        self.oomagent.terminate()

    def online_get(self, entity_key, feature_names):
        with grpc.insecure_channel(self.addr) as channel:
            stub = oomagent_pb2_grpc.OomAgentStub(channel)
            response = stub.OnlineGet(
                oomagent_pb2.OnlineGetRequest(
                    entity_key=entity_key, feature_names=feature_names
                )
            )
        return map_container_to_dict(response.result.map)

    def online_multi_get(self, entity_keys, feature_names):
        with grpc.insecure_channel(self.addr) as channel:
            stub = oomagent_pb2_grpc.OomAgentStub(channel)
            response = stub.OnlineMultiGet(
                oomagent_pb2.OnlineMultiGetRequest(
                    entity_keys=entity_keys, feature_names=feature_names
                )
            )
        return dict(
            {
                entity_key: map_container_to_dict(values.map)
                for entity_key, values in response.result.items()
            }
        )

    def sync(self, revision_id, purge_delay):
        with grpc.insecure_channel(self.addr) as channel:
            stub = oomagent_pb2_grpc.OomAgentStub(channel)
            stub.Sync(
                oomagent_pb2.SyncRequest(
                    revision_id=revision_id,
                    purge_delay=purge_delay,
                )
            )
        return

    def import_(
        self, group_name, description, input_file_path, delimiter, revision=None
    ):
        with grpc.insecure_channel(self.addr) as channel:
            stub = oomagent_pb2_grpc.OomAgentStub(channel)
            response = stub.Import(
                oomagent_pb2.ImportRequest(
                    group_name=group_name,
                    description=description,
                    input_file_path=input_file_path,
                    delimiter=delimiter,
                    revision=revision,
                )
            )
        return response.revision_id

    def join(self, feature_names, input_file_path, output_file_path):
        with grpc.insecure_channel(self.addr) as channel:
            stub = oomagent_pb2_grpc.OomAgentStub(channel)
            stub.Join(
                oomagent_pb2.JoinRequest(
                    feature_names=feature_names,
                    input_file_path=input_file_path,
                    output_file_path=output_file_path,
                )
            )
        return

    def export(self, feature_names, revision_id, output_file_path, limit=None):
        with grpc.insecure_channel(self.addr) as channel:
            stub = oomagent_pb2_grpc.OomAgentStub(channel)
            stub.Export(
                oomagent_pb2.ExportRequest(
                    feature_names=feature_names,
                    revision_id=revision_id,
                    output_file_path=output_file_path,
                    limit=limit,
                )
            )
        return

    def channel_export(self, feature_names, revision_id, limit=None):
        with grpc.insecure_channel(self.addr) as channel:
            stub = oomagent_pb2_grpc.OomAgentStub(channel)
            response_channel = stub.ChannelExport(
                oomagent_pb2.ExportRequest(
                    feature_names=feature_names, revision_id=revision_id, limit=limit
                )
            )
        return response_channel
