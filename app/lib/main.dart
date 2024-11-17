import 'package:flutter/material.dart';
import './candles/types.dart';
import './candles/chart.dart';

void main() {
  runApp(MyApp());
}

// Root of the application
class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Solana App',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: PrivateKeyPage(),
    );
  }
}

// Page to input Solana private key
class PrivateKeyPage extends StatefulWidget {
  @override
  _PrivateKeyPageState createState() => _PrivateKeyPageState();
}

class _PrivateKeyPageState extends State<PrivateKeyPage> {
  final TextEditingController _privateKeyController = TextEditingController();

  void _submitPrivateKey() {
    String privateKey = _privateKeyController.text.trim();
    if (privateKey.isNotEmpty) {
      Navigator.push(
        context,
        MaterialPageRoute(
          builder: (context) => HomePage(privateKey: privateKey),
        ),
      );
    } else {
      // Show an alert if the private key is empty
      showDialog(
        context: context,
        builder: (_) => AlertDialog(
          title: Text('Error'),
          content: Text('Please enter your private key.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: Text('OK'),
            ),
          ],
        ),
      );
    }
  }

  @override
  void dispose() {
    _privateKeyController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: Text('Enter Solana Private Key'),
        ),
        body: Center(
          child: Padding(
            padding: EdgeInsets.all(16.0),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                TextField(
                  controller: _privateKeyController,
                  decoration: InputDecoration(
                    labelText: 'Private Key',
                    border: OutlineInputBorder(),
                  ),
                  obscureText: true, // Hide the input for security
                ),
                SizedBox(height: 20),
                ElevatedButton(
                  onPressed: _submitPrivateKey,
                  child: Text('Submit'),
                ),
              ],
            ),
          ),
        ));
  }
}

// Home page after entering the private key
class HomePage extends StatefulWidget {
  final String privateKey;

  HomePage({required this.privateKey});

  @override
  _HomePageState createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  double _balance = 0.0;
  // Placeholder for chart data
  List<double> _chartData = [];

  @override
  void initState() {
    super.initState();
    _fetchBalance();
    _fetchChartData();
  }

  // Fetch user's balance from your API
  void _fetchBalance() async {
    // TODO: Implement API call to fetch balance
    setState(() {
      _balance = 100.0; // Placeholder value
    });
  }

  // Fetch chart data from your API
  void _fetchChartData() async {
    // TODO: Implement API call to fetch chart data
    setState(() {
      _chartData = [1, 2, 3, 4, 5]; // Placeholder data
    });
  }

  void _buyCoin() {
    // TODO: Implement buy functionality
    print('Buy button pressed');
  }

  void _sellCoin() {
    // TODO: Implement sell functionality
    print('Sell button pressed');
  }

  @override
  Widget build(BuildContext context) {
    // Example coin data
    Coin coinData = Coin(
      id: 'pnut',
      image: 'https://s2.coinmarketcap.com/static/img/coins/64x64/1.png',
      name: 'Bitcoin',
      shortName: 'BTC',
      price: '123456',
      lastPrice: '123456',
      percentage: '-0.5',
      symbol: 'pnut',
      highDay: '567',
      lowDay: '12',
      decimalCurrency: 18,
    );

    return Scaffold(
      appBar: AppBar(
        title: Text('Solana Dashboard'),
      ),
      body: Padding(
        padding: EdgeInsets.all(16.0),
        child: Column(
          children: [
            // Display user's balance
            Text(
              'Balance: \$$_balance',
              style: TextStyle(fontSize: 24),
            ),
            SizedBox(height: 20),
            // Candlestick chart
            Expanded(
              child: CandleChart(
                coinData: coinData,
                inrRate: 77.0,
                intervalSelectedTextColor: Colors.red,
                intervalTextSize: 20,
                intervalUnselectedTextColor: Colors.black,
              ),
            ),
            SizedBox(height: 20),
            // Buy and Sell buttons
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceEvenly,
              children: [
                ElevatedButton(
                  onPressed: _buyCoin,
                  child: Text('Buy'),
                ),
                ElevatedButton(
                  onPressed: _sellCoin,
                  child: Text('Sell'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
