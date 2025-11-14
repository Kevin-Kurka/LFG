import 'package:flutter/foundation.dart';
import 'package:lfg_app/models/order.dart';
import 'package:lfg_app/services/api_client.dart';

class OrderProvider with ChangeNotifier {
  final String? token;
  List<Order> _orders = [];
  bool _isLoading = false;
  String? _error;

  List<Order> get orders => _orders;
  List<Order> get activeOrders => _orders.where((o) => o.isActive).toList();
  bool get isLoading => _isLoading;
  String? get error => _error;

  OrderProvider(this.token);

  Future<void> fetchOrders({String? status}) async {
    if (token == null) return;

    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final apiClient = ApiClient(token: token);
      _orders = await apiClient.getOrders(status: status);
      _isLoading = false;
      notifyListeners();
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<Order> placeOrder({
    required String contractId,
    required String type,
    required int quantity,
    double? limitPrice,
  }) async {
    if (token == null) throw Exception('Not authenticated');

    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final apiClient = ApiClient(token: token);
      final order = await apiClient.placeOrder(
        contractId: contractId,
        type: type,
        quantity: quantity,
        limitPrice: limitPrice,
      );
      await fetchOrders();
      _isLoading = false;
      notifyListeners();
      return order;
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
      rethrow;
    }
  }

  Future<void> cancelOrder(String orderId) async {
    if (token == null) return;

    try {
      final apiClient = ApiClient(token: token);
      await apiClient.cancelOrder(orderId);
      await fetchOrders();
      notifyListeners();
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      rethrow;
    }
  }

  Future<void> refresh() async {
    await fetchOrders();
  }
}
