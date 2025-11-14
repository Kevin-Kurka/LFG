import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import 'package:lfg_app/providers/order_provider.dart';
import 'package:lfg_app/providers/wallet_provider.dart';

class PlaceOrderScreen extends StatefulWidget {
  final String contractId;
  final String marketQuestion;
  final String contractSide;
  final double currentPrice;

  const PlaceOrderScreen({
    super.key,
    required this.contractId,
    required this.marketQuestion,
    required this.contractSide,
    required this.currentPrice,
  });

  @override
  State<PlaceOrderScreen> createState() => _PlaceOrderScreenState();
}

class _PlaceOrderScreenState extends State<PlaceOrderScreen> {
  final _formKey = GlobalKey<FormState>();
  final _quantityController = TextEditingController();
  final _limitPriceController = TextEditingController();

  String _orderType = 'MARKET';
  bool _isPlacingOrder = false;

  @override
  void initState() {
    super.initState();
    _limitPriceController.text = widget.currentPrice.toStringAsFixed(4);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<WalletProvider>(context, listen: false).fetchWallet();
    });
  }

  @override
  void dispose() {
    _quantityController.dispose();
    _limitPriceController.dispose();
    super.dispose();
  }

  double get _estimatedCost {
    final quantity = int.tryParse(_quantityController.text) ?? 0;
    final price = _orderType == 'LIMIT'
        ? (double.tryParse(_limitPriceController.text) ?? widget.currentPrice)
        : widget.currentPrice;
    return quantity * price;
  }

  Future<void> _handlePlaceOrder() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _isPlacingOrder = true;
    });

    try {
      final orderProvider = Provider.of<OrderProvider>(context, listen: false);
      await orderProvider.placeOrder(
        contractId: widget.contractId,
        type: _orderType,
        quantity: int.parse(_quantityController.text),
        limitPrice: _orderType == 'LIMIT'
            ? double.parse(_limitPriceController.text)
            : null,
      );

      if (!mounted) return;

      // Refresh wallet
      await Provider.of<WalletProvider>(context, listen: false).fetchWallet();

      if (!mounted) return;

      Navigator.of(context).pop();
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Order placed successfully!'),
          backgroundColor: Colors.green,
        ),
      );
    } catch (e) {
      if (!mounted) return;

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(e.toString().replaceAll('Exception: ', '')),
          backgroundColor: Colors.red,
        ),
      );
    } finally {
      if (mounted) {
        setState(() {
          _isPlacingOrder = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Place Order'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        widget.marketQuestion,
                        style: Theme.of(context).textTheme.titleMedium?.copyWith(
                              fontWeight: FontWeight.bold,
                            ),
                      ),
                      const SizedBox(height: 12),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            widget.contractSide,
                            style: Theme.of(context).textTheme.titleLarge?.copyWith(
                                  color: Theme.of(context).colorScheme.primary,
                                  fontWeight: FontWeight.bold,
                                ),
                          ),
                          Text(
                            '${widget.currentPrice.toStringAsFixed(4)} credits',
                            style: Theme.of(context).textTheme.bodyLarge,
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 24),
              Consumer<WalletProvider>(
                builder: (context, wallet, child) {
                  return Card(
                    color: Theme.of(context).colorScheme.primaryContainer,
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            'Available Balance',
                            style: Theme.of(context).textTheme.titleMedium,
                          ),
                          Text(
                            '${wallet.wallet?.balance.toStringAsFixed(2) ?? '0.00'} credits',
                            style: Theme.of(context).textTheme.titleMedium?.copyWith(
                                  fontWeight: FontWeight.bold,
                                ),
                          ),
                        ],
                      ),
                    ),
                  );
                },
              ),
              const SizedBox(height: 24),
              Text(
                'Order Type',
                style: Theme.of(context).textTheme.titleMedium,
              ),
              const SizedBox(height: 12),
              SegmentedButton<String>(
                segments: const [
                  ButtonSegment(
                    value: 'MARKET',
                    label: Text('Market'),
                  ),
                  ButtonSegment(
                    value: 'LIMIT',
                    label: Text('Limit'),
                  ),
                ],
                selected: {_orderType},
                onSelectionChanged: (Set<String> newSelection) {
                  setState(() {
                    _orderType = newSelection.first;
                  });
                },
              ),
              const SizedBox(height: 24),
              TextFormField(
                controller: _quantityController,
                keyboardType: TextInputType.number,
                inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                decoration: const InputDecoration(
                  labelText: 'Quantity',
                  prefixIcon: Icon(Icons.numbers),
                  border: OutlineInputBorder(),
                  helperText: 'Number of contracts to buy',
                ),
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter quantity';
                  }
                  final quantity = int.tryParse(value);
                  if (quantity == null || quantity <= 0) {
                    return 'Quantity must be greater than 0';
                  }
                  return null;
                },
                onChanged: (_) => setState(() {}),
              ),
              const SizedBox(height: 20),
              if (_orderType == 'LIMIT') ...[
                TextFormField(
                  controller: _limitPriceController,
                  keyboardType: const TextInputType.numberWithOptions(decimal: true),
                  decoration: const InputDecoration(
                    labelText: 'Limit Price (credits)',
                    prefixIcon: Icon(Icons.attach_money),
                    border: OutlineInputBorder(),
                    helperText: 'Maximum price per contract',
                  ),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Please enter limit price';
                    }
                    final price = double.tryParse(value);
                    if (price == null || price <= 0) {
                      return 'Price must be greater than 0';
                    }
                    return null;
                  },
                  onChanged: (_) => setState(() {}),
                ),
                const SizedBox(height: 20),
              ],
              Card(
                color: Theme.of(context).colorScheme.secondaryContainer,
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const Text('Estimated Cost'),
                          Text(
                            '${_estimatedCost.toStringAsFixed(2)} credits',
                            style: const TextStyle(
                              fontWeight: FontWeight.bold,
                              fontSize: 18,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 30),
              FilledButton(
                onPressed: _isPlacingOrder ? null : _handlePlaceOrder,
                style: FilledButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 16),
                ),
                child: _isPlacingOrder
                    ? const SizedBox(
                        height: 20,
                        width: 20,
                        child: CircularProgressIndicator(strokeWidth: 2),
                      )
                    : const Text('Place Order', style: TextStyle(fontSize: 16)),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
