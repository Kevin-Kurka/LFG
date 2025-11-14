import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:lfg_app/models/user.dart';
import 'package:lfg_app/models/market.dart';
import 'package:lfg_app/models/order.dart';
import 'package:lfg_app/models/wallet.dart';

class ApiClient {
  // Change this to your API Gateway URL
  static const String baseUrl = 'http://localhost:8000';

  final String? token;

  ApiClient({this.token});

  Map<String, String> get _headers => {
        'Content-Type': 'application/json',
        if (token != null) 'Authorization': 'Bearer $token',
      };

  // Auth endpoints
  Future<Map<String, dynamic>> register(String email, String password, String displayName) async {
    final response = await http.post(
      Uri.parse('$baseUrl/register'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'email': email,
        'password': password,
        'display_name': displayName,
      }),
    );

    if (response.statusCode == 201) {
      return jsonDecode(response.body);
    } else {
      throw Exception(jsonDecode(response.body)['error'] ?? 'Registration failed');
    }
  }

  Future<Map<String, dynamic>> login(String email, String password) async {
    final response = await http.post(
      Uri.parse('$baseUrl/login'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'email': email,
        'password': password,
      }),
    );

    if (response.statusCode == 200) {
      return jsonDecode(response.body);
    } else {
      throw Exception(jsonDecode(response.body)['error'] ?? 'Login failed');
    }
  }

  Future<User> getProfile() async {
    final response = await http.get(
      Uri.parse('$baseUrl/profile'),
      headers: _headers,
    );

    if (response.statusCode == 200) {
      return User.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to fetch profile');
    }
  }

  // Market endpoints
  Future<List<Market>> getMarkets({
    String? status,
    String? search,
    int page = 1,
    int pageSize = 20,
  }) async {
    final queryParams = {
      if (status != null) 'status': status,
      if (search != null) 'search': search,
      'page': page.toString(),
      'page_size': pageSize.toString(),
    };

    final uri = Uri.parse('$baseUrl/markets').replace(queryParameters: queryParams);
    final response = await http.get(uri, headers: _headers);

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return (data['markets'] as List).map((m) => Market.fromJson(m)).toList();
    } else {
      throw Exception('Failed to fetch markets');
    }
  }

  Future<Market> getMarket(String marketId) async {
    final response = await http.get(
      Uri.parse('$baseUrl/markets/$marketId'),
      headers: _headers,
    );

    if (response.statusCode == 200) {
      return Market.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to fetch market');
    }
  }

  // Wallet endpoints
  Future<Wallet> getWallet() async {
    final response = await http.get(
      Uri.parse('$baseUrl/wallet'),
      headers: _headers,
    );

    if (response.statusCode == 200) {
      return Wallet.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to fetch wallet');
    }
  }

  Future<List<WalletTransaction>> getWalletTransactions({int page = 1, int pageSize = 20}) async {
    final uri = Uri.parse('$baseUrl/wallet/transactions').replace(
      queryParameters: {
        'page': page.toString(),
        'page_size': pageSize.toString(),
      },
    );

    final response = await http.get(uri, headers: _headers);

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return (data['transactions'] as List).map((t) => WalletTransaction.fromJson(t)).toList();
    } else {
      throw Exception('Failed to fetch transactions');
    }
  }

  // Order endpoints
  Future<Order> placeOrder({
    required String contractId,
    required String type,
    required int quantity,
    double? limitPrice,
  }) async {
    final response = await http.post(
      Uri.parse('$baseUrl/orders'),
      headers: _headers,
      body: jsonEncode({
        'contract_id': contractId,
        'type': type,
        'quantity': quantity,
        if (limitPrice != null) 'limit_price_credits': limitPrice,
      }),
    );

    if (response.statusCode == 201) {
      return Order.fromJson(jsonDecode(response.body));
    } else {
      throw Exception(jsonDecode(response.body)['error'] ?? 'Failed to place order');
    }
  }

  Future<void> cancelOrder(String orderId) async {
    final response = await http.post(
      Uri.parse('$baseUrl/orders/$orderId/cancel'),
      headers: _headers,
    );

    if (response.statusCode != 200) {
      throw Exception(jsonDecode(response.body)['error'] ?? 'Failed to cancel order');
    }
  }

  Future<List<Order>> getOrders({String? status, int page = 1, int pageSize = 20}) async {
    final queryParams = {
      if (status != null) 'status': status,
      'page': page.toString(),
      'page_size': pageSize.toString(),
    };

    final uri = Uri.parse('$baseUrl/orders').replace(queryParameters: queryParams);
    final response = await http.get(uri, headers: _headers);

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return (data['orders'] as List).map((o) => Order.fromJson(o)).toList();
    } else {
      throw Exception('Failed to fetch orders');
    }
  }

  // Credit exchange endpoints
  Future<void> buyCredits({
    required String cryptoType,
    required double cryptoAmount,
  }) async {
    final response = await http.post(
      Uri.parse('$baseUrl/credit-exchange/buy'),
      headers: _headers,
      body: jsonEncode({
        'crypto_type': cryptoType,
        'crypto_amount': cryptoAmount,
      }),
    );

    if (response.statusCode != 200) {
      throw Exception(jsonDecode(response.body)['error'] ?? 'Failed to buy credits');
    }
  }
}
