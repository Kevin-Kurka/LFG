class Wallet {
  final String id;
  final String userId;
  final double balance;
  final DateTime updatedAt;

  Wallet({
    required this.id,
    required this.userId,
    required this.balance,
    required this.updatedAt,
  });

  factory Wallet.fromJson(Map<String, dynamic> json) {
    return Wallet(
      id: json['id'],
      userId: json['user_id'],
      balance: (json['balance_credits'] ?? 0).toDouble(),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }
}

class WalletTransaction {
  final String id;
  final String walletId;
  final double amount;
  final String type;
  final String description;
  final DateTime createdAt;

  WalletTransaction({
    required this.id,
    required this.walletId,
    required this.amount,
    required this.type,
    required this.description,
    required this.createdAt,
  });

  factory WalletTransaction.fromJson(Map<String, dynamic> json) {
    return WalletTransaction(
      id: json['id'],
      walletId: json['wallet_id'],
      amount: (json['amount_credits']).toDouble(),
      type: json['transaction_type'],
      description: json['description'] ?? '',
      createdAt: DateTime.parse(json['created_at']),
    );
  }
}
