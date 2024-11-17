class Candle {
  /// DateTime for the candle
  final DateTime date;

  /// The highest price during this candle lifetime
  /// It if always more than low, open and close
  final double high;

  /// The lowest price during this candle lifetime
  /// It if always less than high, open and close
  final double low;

  /// Price at the beginning of the period
  final double open;

  /// Price at the end of the period
  final double close;

  /// Volume is the number of shares of a
  /// security traded during a given period of time.
  final double volume;

  bool get isBull => open <= close;

  Candle({
    required this.date,
    required this.high,
    required this.low,
    required this.open,
    required this.close,
    required this.volume,
  });

  Candle.fromJson(List<dynamic> json)
      : date = DateTime.fromMillisecondsSinceEpoch(json[0]),
        open = json[1].toDouble(),
        high = json[2].toDouble(),
        low = json[3].toDouble(),
        close = json[4].toDouble(),
        volume = json[5].toDouble();
}

/// Coin model which holds a single coin data.
/// It contains 14 required variables that hold a single coin data:
/// coinID, coinImage, coinName, coinShortName, coinPrice,
/// coinLastPrice, coinSymbol, coinPairWith, coinHighDay, coinLowDay,
/// coinDecimalPair, coinDecimalCurrency and coinListed.
///

class Coin {
  String id;
  String image;
  String name;
  String shortName;
  String price;
  String lastPrice;
  String percentage;
  String symbol;
  String highDay;
  String lowDay;
  int decimalCurrency;

  Coin({
    required this.id,
    required this.image,
    required this.name,
    required this.shortName,
    required this.price,
    required this.lastPrice,
    required this.percentage,
    required this.symbol,
    required this.highDay,
    required this.lowDay,
    required this.decimalCurrency,
  });

  @override
  String toString() {
    return 'Coin{coinID: $id, coinImage: $image, coinName: $name, coinShortName: $shortName, coinPrice: $price, coinLastPrice: $lastPrice, coinPercentage: $percentage, coinSymbol: $symbol, coinHighDay: $highDay, coinLowDay: $lowDay, coinDecimalCurrency: $decimalCurrency}';
  }
}
