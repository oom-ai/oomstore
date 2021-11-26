import sys
import grpc
import logging
from .codegen import oomd_pb2
from .codegen import oomd_pb2_grpc
import time
from pathlib import Path
from subprocess import Popen

# Convert google.protobuf.pyext._message.MessageMapContainer object to Python dictionary
def map_container_to_dict(map_container):
    return dict({k: getattr(v, v.WhichOneof("kind")) for k, v in map_container.items()})


class Client(object):
    def __init__(self, port: int, config_path: str):
        try:
            self.oomd = Popen(["oomd", "--config", config_path])
        except Exception as e:
            logging.error(e)
            sys.exit(1)

        self.addr = "127.0.0.1:%d" % port

        # wait for oomd to start
        time.sleep(2)

    def __del__(self):
        self.oomd.terminate()

    def online_get(self, entity_key, feature_names):
        with grpc.insecure_channel(self.addr) as channel:
            stub = codegen.oomd_pb2_grpc.OomDStub(channel)
            response = stub.OnlineGet(
                codegen.oomd_pb2.OnlineGetRequest(
                    entity_key=entity_key, feature_names=feature_names
                )
            )
        return map_container_to_dict(response.result.map)

    def online_multi_get(self, entity_keys, feature_names):
        with grpc.insecure_channel(self.addr) as channel:
            stub = codegen.oomd_pb2_grpc.OomDStub(channel)
            response = stub.OnlineMultiGet(
                codegen.oomd_pb2.OnlineMultiGetRequest(
                    entity_keys=entity_keys, feature_names=feature_names
                )
            )
        return dict(
            {
                entity_key: map_container_to_dict(values.map)
                for entity_key, values in response.result.items()
            }
        )

    def sync(self, revision_id):
        with grpc.insecure_channel(self.addr) as channel:
            stub = codegen.oomd_pb2_grpc.OomDStub(channel)
            stub.Sync(codegen.oomd_pb2.SyncRequest(revision_id=revision_id))
        return

    def import_(
        self, group_name, description, input_file_path, delimiter, revision=None
    ):
        with grpc.insecure_channel(self.addr) as channel:
            stub = codegen.oomd_pb2_grpc.OomDStub(channel)
            response = stub.Import(
                codegen.oomd_pb2.ImportRequest(
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
            stub = codegen.oomd_pb2_grpc.OomDStub(channel)
            stub.Join(
                codegen.oomd_pb2.JoinRequest(
                    feature_names=feature_names,
                    input_file_path=input_file_path,
                    output_file_path=output_file_path,
                )
            )
        return

    def export(self, feature_names, revision_id, output_file_path, limit=None):
        with grpc.insecure_channel(self.addr) as channel:
            stub = codegen.oomd_pb2_grpc.OomDStub(channel)
            stub.Export(
                codegen.oomd_pb2.ExportRequest(
                    feature_names=feature_names,
                    revision_id=revision_id,
                    output_file_path=output_file_path,
                    limit=limit
                )
            )
        return

    def channel_export(self, feature_names, revision_id, limit=None):
        with grpc.insecure_channel(self.addr) as channel:
            stub = codegen.oomd_pb2_grpc.OomDStub(channel)
            response_channel = stub.ChannelExport(
                codegen.oomd_pb2.ExportRequest(
                    feature_names=feature_names,
                    revision_id=revision_id,
                    limit=limit
                )
            )
        return response_channel


if __name__ == "__main__":
    config_path = "%s/.config/oomstore/config.yaml" % str(Path.home())
    client = Client(50051, config_path)
    revision_id1 = client.import_(
        group_name="account",
        description="sample account data",
        input_file_path="/tmp/account.csv",
        delimiter=",",
    )
    revision_id2 = client.import_(
        group_name="transaction_stats",
        description="sample transaction stat data",
        input_file_path="/tmp/transaction_stats.csv",
        delimiter=",",
    )
    client.sync(revision_id1)
    client.sync(revision_id2)
    print(
        client.online_get(
            entity_key="1006",
            feature_names=[
                "state",
                "credit_score",
                "account_age_days",
                "has_2fa_installed",
                "transaction_count_7d",
                "transaction_count_30d",
            ],
        )
    )
    print(
        client.online_multi_get(
            entity_keys=["1006", "1007"],
            feature_names=[
                "state",
                "credit_score",
                "account_age_days",
                "has_2fa_installed",
                "transaction_count_7d",
                "transaction_count_30d",
            ],
        )
    )
    client.join(
        feature_names=[
            "state",
            "credit_score",
            "account_age_days",
            "has_2fa_installed",
            "transaction_count_7d",
            "transaction_count_30d",
        ],
        input_file_path="/tmp/label.csv",
        output_file_path="/tmp/joined.csv",
    )
