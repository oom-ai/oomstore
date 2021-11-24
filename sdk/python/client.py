import grpc
import codegen.oomd_pb2
import codegen.oomd_pb2_grpc
import time

class Client(object):
    def __init__(self):
        pass

    def online_get(self, entity_key, feature_names):
        with grpc.insecure_channel('localhost:50051') as channel:
            stub = codegen.oomd_pb2_grpc.OomDStub(channel)
            response = stub.OnlineGet(codegen.oomd_pb2.OnlineGetRequest(entity_key=entity_key, feature_names=feature_names))
        return response.result

    def online_multi_get(self, entity_keys, feature_names):
        with grpc.insecure_channel('localhost:50051') as channel:
            stub = codegen.oomd_pb2_grpc.OomDStub(channel)
            response = stub.OnlineMultiGet(codegen.oomd_pb2.OnlineMultiGetRequest(entity_keys=entity_keys, feature_names=feature_names))
        return response.result

    def sync(self, revision_id):
        with grpc.insecure_channel('localhost:50051') as channel:
            stub = codegen.oomd_pb2_grpc.OomDStub(channel)
            stub.Sync(codegen.oomd_pb2.SyncRequest(revision_id=revision_id))
        return

    def import_by_file(self, group_name, description, input_file_path, delimiter, revision=None):
        with grpc.insecure_channel('localhost:50051') as channel:
            stub = codegen.oomd_pb2_grpc.OomDStub(channel)
            response = stub.ImportByFile(codegen.oomd_pb2.ImportByFileRequest(group_name=group_name,description=description,input_file_path=input_file_path,delimiter=delimiter,revision=revision))
        return response.revision_id

    def join_by_file(self, feature_names, input_file_path, output_file_path):
        with grpc.insecure_channel('localhost:50051') as channel:
            stub = codegen.oomd_pb2_grpc.OomDStub(channel)
            stub.JoinByFile(codegen.oomd_pb2.JoinByFileRequest(feature_names=feature_names,input_file_path=input_file_path,output_file_path=output_file_path))
        return

if __name__ == "__main__":
    client = Client()
    revision_id1 = client.import_by_file(group_name='account', description='sample account data', input_file_path='/tmp/account.csv', delimiter=',')
    revision_id2 = client.import_by_file(group_name='transaction_stats', description='sample transaction stat data', input_file_path='/tmp/transaction_stats.csv', delimiter=',')
    time.sleep(10)
    client.sync(revision_id1)
    client.sync(revision_id2)
    print(client.online_get(entity_key="1006", feature_names=["state", "credit_score", "account_age_days", "has_2fa_installed", "transaction_count_7d", "transaction_count_30d"]))
    client.join_by_file(feature_names=["state", "credit_score", "account_age_days", "has_2fa_installed", "transaction_count_7d", "transaction_count_30d"], input_file_path='/tmp/label.csv', output_file_path='/tmp/joined.csv')
