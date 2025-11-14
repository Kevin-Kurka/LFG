import 'package:flutter/foundation.dart';
import 'package:lfg_app/models/wallet.dart';
import 'package:lfg_app/services/api_client.dart';

class WalletProvider with ChangeNotifier {
  final String? token;
  Wallet? _wallet;
  List<WalletTransaction> _transactions = [];
  bool _isLoading = false;
  String? _error;

  Wallet? get wallet => _wallet;
  List<WalletTransaction> get transactions => _transactions;
  bool get isLoading => _isLoading;
  String? get error => _error;

  WalletProvider(this.token);

  Future<void> fetchWallet() async {
    if (token == null) return;

    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final apiClient = ApiClient(token: token);
      _wallet = await apiClient.getWallet();
      _isLoading = false;
      notifyListeners();
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> fetchTransactions() async {
    if (token == null) return;

    try {
      final apiClient = ApiClient(token: token);
      _transactions = await apiClient.getWalletTransactions();
      notifyListeners();
    } catch (e) {
      _error = e.toString();
      notifyListeners();
    }
  }

  Future<void> buyCredits(String cryptoType, double cryptoAmount) async {
    if (token == null) return;

    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final apiClient = ApiClient(token: token);
      await apiClient.buyCredits(
        cryptoType: cryptoType,
        cryptoAmount: cryptoAmount,
      );
      await fetchWallet();
      await fetchTransactions();
      _isLoading = false;
      notifyListeners();
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
      rethrow;
    }
  }

  Future<void> refresh() async {
    await fetchWallet();
    await fetchTransactions();
  }
}
