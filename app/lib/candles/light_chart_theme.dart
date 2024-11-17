import 'package:k_chart/chart_style.dart';
import 'package:flutter/material.dart' show Color;

// Custom ChartColors with light mode
// Also flips green and red candles (original library was Chinese and the colors are flipped LOL)
class CustomChartColors extends ChartColors {
  @override
  List<Color> bgColor = [Color(0x00FFF4FC), Color(0x00FFF4FC)];

  @override
  Color kLineColor = Color(0xff3165A4);

  @override
  Color lineFillColor = Color(0x553165A4);

  @override
  Color lineFillInsideColor = Color(0xffffffff);

  @override
  Color ma5Color = Color(0xff785D1B);

  @override
  Color ma10Color = Color(0xff3F8179);

  @override
  Color ma30Color = Color(0xff68499A);

  @override
  Color upColor = Color(0xffAB3A46);

  @override
  Color dnColor = Color(0xff3A885B);

  @override
  Color volColor = Color(0xff7452E2);

  @override
  Color macdColor = Color(0xff7452E2);

  @override
  Color difColor = Color(0xff785D1B);

  @override
  Color deaColor = Color(0xff3F8179);

  @override
  Color kColor = Color(0xff785D1B);

  @override
  Color dColor = Color(0xff3F8179);

  @override
  Color jColor = Color(0xff68499A);

  @override
  Color rsiColor = Color(0xff785D1B);

  @override
  Color defaultTextColor = Color(0xff404040);

  @override
  Color nowPriceUpColor = Color(0xff3A885B);

  @override
  Color nowPriceDnColor = Color(0xffAB3A46);

  @override
  Color nowPriceTextColor = Color(0xff000000);

  @override
  Color depthBuyColor = Color(0xff4B9077);

  @override
  Color depthSellColor = Color(0xffAB3A46);

  @override
  Color selectBorderColor = Color(0xffBCC6CF);

  @override
  Color selectFillColor = Color(0xffE8EEF3);

  @override
  Color gridColor = Color(0xffD6D6D6);

  @override
  Color infoWindowNormalColor = Color(0xff000000);

  @override
  Color infoWindowTitleColor = Color(0xff000000);

  @override
  Color infoWindowUpColor = Color(0xff008800);

  @override
  Color infoWindowDnColor = Color(0xff880000);

  @override
  Color hCrossColor = Color(0xff000000);

  @override
  Color vCrossColor = Color(0x1E000000);

  @override
  Color crossTextColor = Color(0xff000000);

  @override
  Color maxColor = Color(0xff000000);

  @override
  Color minColor = Color(0xff000000);

  @override
  Color getMAColor(int index) {
    switch (index % 3) {
      case 1:
        return ma10Color;
      case 2:
        return ma30Color;
      default:
        return ma5Color;
    }
  }
}

// Custom ChartStyle with overrides
class CustomChartStyle extends ChartStyle {
  @override
  double topPadding = 40.0;

  @override
  double bottomPadding = 40.0;

  @override
  double childPadding = 15.0;

  @override
  double pointWidth = 12.0;

  @override
  double candleWidth = 9.0;

  @override
  double candleLineWidth = 2.0;

  @override
  double volWidth = 9.0;

  @override
  double macdWidth = 4.0;

  @override
  double vCrossWidth = 9.0;

  @override
  double hCrossWidth = 0.7;

  @override
  double nowPriceLineLength = 2.0;

  @override
  double nowPriceLineSpan = 2.0;

  @override
  double nowPriceLineWidth = 1.5;

  @override
  int gridRows = 5;

  @override
  int gridColumns = 5;
}

final ChartStyle chartStyle = CustomChartStyle();
final ChartColors chartColors = CustomChartColors();
