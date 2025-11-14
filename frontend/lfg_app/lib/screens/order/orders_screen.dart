import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:lfg_app/providers/order_provider.dart';
import 'package:intl/intl.dart';

class OrdersScreen extends StatefulWidget {
  const OrdersScreen({super.key});

  @override
  State<OrdersScreen> createState() => _OrdersScreenState();
}

class _OrdersScreenState extends State<OrdersScreen> {
  String? _statusFilter;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<OrderProvider>(context, listen: false).fetchOrders();
    });
  }

  Future<void> _handleCancelOrder(String orderId) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Cancel Order'),
        content: const Text('Are you sure you want to cancel this order?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('No'),
          ),
          FilledButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: const Text('Yes, Cancel'),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      try {
        await Provider.of<OrderProvider>(context, listen: false).cancelOrder(orderId);
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('Order cancelled successfully'),
              backgroundColor: Colors.green,
            ),
          );
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(e.toString().replaceAll('Exception: ', '')),
              backgroundColor: Colors.red,
            ),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Orders'),
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(60),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: Row(
                children: [
                  FilterChip(
                    label: const Text('All'),
                    selected: _statusFilter == null,
                    onSelected: (selected) {
                      setState(() {
                        _statusFilter = null;
                      });
                      Provider.of<OrderProvider>(context, listen: false).fetchOrders();
                    },
                  ),
                  const SizedBox(width: 8),
                  FilterChip(
                    label: const Text('Active'),
                    selected: _statusFilter == 'ACTIVE',
                    onSelected: (selected) {
                      setState(() {
                        _statusFilter = selected ? 'ACTIVE' : null;
                      });
                      Provider.of<OrderProvider>(context, listen: false).fetchOrders(status: _statusFilter);
                    },
                  ),
                  const SizedBox(width: 8),
                  FilterChip(
                    label: const Text('Filled'),
                    selected: _statusFilter == 'FILLED',
                    onSelected: (selected) {
                      setState(() {
                        _statusFilter = selected ? 'FILLED' : null;
                      });
                      Provider.of<OrderProvider>(context, listen: false).fetchOrders(status: _statusFilter);
                    },
                  ),
                  const SizedBox(width: 8),
                  FilterChip(
                    label: const Text('Cancelled'),
                    selected: _statusFilter == 'CANCELLED',
                    onSelected: (selected) {
                      setState(() {
                        _statusFilter = selected ? 'CANCELLED' : null;
                      });
                      Provider.of<OrderProvider>(context, listen: false).fetchOrders(status: _statusFilter);
                    },
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
      body: Consumer<OrderProvider>(
        builder: (context, provider, child) {
          if (provider.isLoading && provider.orders.isEmpty) {
            return const Center(child: CircularProgressIndicator());
          }

          if (provider.error != null) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.error_outline, size: 48, color: Colors.red),
                  const SizedBox(height: 16),
                  Text('Error: ${provider.error}'),
                  const SizedBox(height: 16),
                  FilledButton(
                    onPressed: () => provider.fetchOrders(),
                    child: const Text('Retry'),
                  ),
                ],
              ),
            );
          }

          if (provider.orders.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.list_alt,
                    size: 64,
                    color: Colors.grey[400],
                  ),
                  const SizedBox(height: 16),
                  Text(
                    'No orders yet',
                    style: Theme.of(context).textTheme.titleLarge,
                  ),
                  const SizedBox(height: 8),
                  const Text('Place your first order to get started'),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: () => provider.refresh(),
            child: ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: provider.orders.length,
              itemBuilder: (context, index) {
                final order = provider.orders[index];
                final fillPercentage = order.quantity > 0
                    ? (order.quantityFilled / order.quantity * 100).toStringAsFixed(0)
                    : '0';

                Color statusColor;
                switch (order.status) {
                  case 'FILLED':
                    statusColor = Colors.green;
                    break;
                  case 'ACTIVE':
                  case 'PARTIALLY_FILLED':
                    statusColor = Colors.blue;
                    break;
                  case 'CANCELLED':
                    statusColor = Colors.orange;
                    break;
                  case 'REJECTED':
                    statusColor = Colors.red;
                    break;
                  default:
                    statusColor = Colors.grey;
                }

                return Card(
                  margin: const EdgeInsets.only(bottom: 12),
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Chip(
                              label: Text(order.type),
                              materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                            ),
                            Chip(
                              label: Text(order.status),
                              backgroundColor: statusColor.withOpacity(0.2),
                              materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                            ),
                          ],
                        ),
                        const SizedBox(height: 12),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            const Text('Quantity:'),
                            Text(
                              '${order.quantityFilled} / ${order.quantity}',
                              style: const TextStyle(fontWeight: FontWeight.bold),
                            ),
                          ],
                        ),
                        const SizedBox(height: 8),
                        LinearProgressIndicator(
                          value: order.quantity > 0 ? order.quantityFilled / order.quantity : 0,
                          backgroundColor: Colors.grey[200],
                          valueColor: AlwaysStoppedAnimation<Color>(statusColor),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          '$fillPercentage% filled',
                          style: Theme.of(context).textTheme.bodySmall,
                        ),
                        if (order.limitPrice != null) ...[
                          const SizedBox(height: 8),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              const Text('Limit Price:'),
                              Text(
                                '${order.limitPrice!.toStringAsFixed(4)} credits',
                                style: const TextStyle(fontWeight: FontWeight.bold),
                              ),
                            ],
                          ),
                        ],
                        const SizedBox(height: 8),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              'Created: ${DateFormat.MMMd().add_jm().format(order.createdAt)}',
                              style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                    color: Colors.grey,
                                  ),
                            ),
                          ],
                        ),
                        if (order.isActive) ...[
                          const SizedBox(height: 12),
                          SizedBox(
                            width: double.infinity,
                            child: OutlinedButton.icon(
                              onPressed: () => _handleCancelOrder(order.id),
                              icon: const Icon(Icons.cancel),
                              label: const Text('Cancel Order'),
                            ),
                          ),
                        ],
                      ],
                    ),
                  ),
                );
              },
            ),
          );
        },
      ),
    );
  }
}
