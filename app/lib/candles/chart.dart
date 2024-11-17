import 'package:flutter/cupertino.dart';
import 'package:get/get.dart';
import 'package:k_chart/chart_style.dart';
import 'package:k_chart/k_chart_widget.dart';

import './types.dart';
import '../widgets/interval.dart';
import './chart_controller.dart';

class CandleChart extends StatelessWidget {
  final Coin coinData;
  final double inrRate;
  final Color backgroundColor;
  final Function()? onSecondaryTap;
  final bool isLine;
  final bool hideGrid;
  final bool hideVolume;
  final bool showNowPrice;
  final bool isTrendLine;
  final bool isTapShowInfoDialog;
  final bool materialInfoDialog;
  final bool showInfoDialog;
  final Color? intervalSelectedTextColor;
  final Color? intervalUnselectedTextColor;
  final double? intervalTextSize;
  final MainAxisAlignment? intervalAlignment;

  CandleChart({
    Key? key,
    required this.coinData,
    required this.inrRate,
    this.backgroundColor = const Color(0xffffffff),
    this.onSecondaryTap,
    this.isLine = false,
    this.hideGrid = false,
    this.hideVolume = true,
    this.showNowPrice = true,
    this.isTrendLine = false,
    this.isTapShowInfoDialog = true,
    this.materialInfoDialog = true,
    this.showInfoDialog = true,
    this.intervalSelectedTextColor,
    this.intervalUnselectedTextColor,
    this.intervalTextSize,
    this.intervalAlignment,
  }) : super(key: key);

  final ChartStyle chartStyle = ChartStyle();
  final ChartColors chartColors = ChartColors();

  @override
  Widget build(BuildContext context) {
    ChartController.to.getCandles(
      coinData: coinData,
      interval: '1',
    );
    ChartController.to.inrRate = inrRate;

    return GetBuilder<ChartController>(
      builder: (_) {
        return Column(
          children: [
            ChartController.to.kChartCandles.isEmpty
                ? const Center(
                    child: CupertinoActivityIndicator(),
                  )
                : Flexible(
                    child: Container(
                      color: backgroundColor,
                      child: KChartWidget(
                        ChartController.to.kChartCandles,
                        chartStyle,
                        chartColors,
                        isLine: isLine,
                        onSecondaryTap: onSecondaryTap,
                        mainState: MainState.NONE,
                        volHidden: hideVolume,
                        fixedLength: 2,
                        timeFormat: TimeFormat.YEAR_MONTH_DAY,
                        hideGrid: hideGrid,
                        showNowPrice: showNowPrice,
                        isTrendLine: isTrendLine,
                        isTapShowInfoDialog: isTapShowInfoDialog,
                        materialInfoDialog: materialInfoDialog,
                        showInfoDialog: showInfoDialog,
                        flingCurve: Curves.bounceInOut,
                      ),
                    ),
                  ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8.0),
              child: Row(
                mainAxisAlignment:
                    intervalAlignment ?? MainAxisAlignment.spaceAround,
                children: [
                  intervalButton(
                    title: '15min',
                    coinData: coinData,
                    intervalSelectedTextColor: intervalSelectedTextColor,
                    intervalUnselectedTextColor: intervalUnselectedTextColor,
                    intervalTextSize: intervalTextSize,
                  ),
                  intervalButton(
                    title: '1h',
                    coinData: coinData,
                    intervalSelectedTextColor: intervalSelectedTextColor,
                    intervalUnselectedTextColor: intervalUnselectedTextColor,
                    intervalTextSize: intervalTextSize,
                  ),
                  intervalButton(
                    title: '4h',
                    coinData: coinData,
                    intervalSelectedTextColor: intervalSelectedTextColor,
                    intervalUnselectedTextColor: intervalUnselectedTextColor,
                    intervalTextSize: intervalTextSize,
                  ),
                  intervalButton(
                    title: '1d',
                    coinData: coinData,
                    intervalSelectedTextColor: intervalSelectedTextColor,
                    intervalUnselectedTextColor: intervalUnselectedTextColor,
                    intervalTextSize: intervalTextSize,
                  )
                ],
              ),
            ),
          ],
        );
      },
    );
  }
}
