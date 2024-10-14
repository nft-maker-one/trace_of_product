pragma circom 2.0.0;

include "circomlib/circuits/poseidon.circom";
include "circomlib/circuits/comparators.circom";
include "circomlib/circuits/bitify.circom";

template MetaDataCircuit() {
    // 输入信号
    signal input EggplantId;
    signal input ProductHeight;
    signal input ProductHash[4];
    signal input TransportHeight;
    signal input TransportHash[4];
    signal input ProcessHeight;
    signal input ProcessHash[4];
    signal input StorageHeight;
    signal input StorageHash[4];
    signal input SellHeight;
    signal input SellHash[4];

    // 输出信号
    signal output valid;

    // 验证高度的顺序
    component ltProduct = LessThan(252);
    ltProduct.in[0] <== ProductHeight;
    ltProduct.in[1] <== TransportHeight;

    component ltTransport = LessThan(252);
    ltTransport.in[0] <== TransportHeight;
    ltTransport.in[1] <== ProcessHeight;

    component ltProcess = LessThan(252);
    ltProcess.in[0] <== ProcessHeight;
    ltProcess.in[1] <== StorageHeight;

    component ltStorage = LessThan(252);
    ltStorage.in[0] <== StorageHeight;
    ltStorage.in[1] <== SellHeight;

    // 验证哈希值
    component hashCheck1 = Poseidon(4);
    component hashCheck2 = Poseidon(4);
    component hashCheck3 = Poseidon(4);
    component hashCheck4 = Poseidon(4);
    component hashCheck5 = Poseidon(4);

    for (var i = 0; i < 4; i++) {
        hashCheck1.inputs[i] <== ProductHash[i];
        hashCheck2.inputs[i] <== TransportHash[i];
        hashCheck3.inputs[i] <== ProcessHash[i];
        hashCheck4.inputs[i] <== StorageHash[i];
        hashCheck5.inputs[i] <== SellHash[i];
    }

    // 使用 AND 门来组合所有条件
    component and1 = AND();
    component and2 = AND();
    component and3 = AND();
    component and4 = AND();
    component and5 = AND();
    component and6 = AND();
    component and7 = AND();
    component and8 = AND();

    and1.a <== ltProduct.out;
    and1.b <== ltTransport.out;

    and2.a <== and1.out;
    and2.b <== ltProcess.out;

    and3.a <== and2.out;
    and3.b <== ltStorage.out;

    and4.a <== and3.out;
    and4.b <== hashCheck1.out;

    and5.a <== and4.out;
    and5.b <== hashCheck2.out;

    and6.a <== and5.out;
    and6.b <== hashCheck3.out;

    and7.a <== and6.out;
    and7.b <== hashCheck4.out;

    and8.a <== and7.out;
    and8.b <== hashCheck5.out;

    valid <== and8.out;
}

template AND() {
    signal input a;
    signal input b;
    signal output out;

    out <== a * b;
}

component main = MetaDataCircuit();