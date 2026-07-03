#ifndef BASECONVERTERTOOL_H
#define BASECONVERTERTOOL_H

#include <QWidget>

class QLineEdit;
class QLabel;

class BaseConverterTool : public QWidget {
    Q_OBJECT

  public:
    explicit BaseConverterTool(QWidget *parent = nullptr);
    ~BaseConverterTool();

  private slots:
    void onBinaryTextChanged(const QString &text);
    void onOctalTextChanged(const QString &text);
    void onDecimalTextChanged(const QString &text);
    void onHexadecimalTextChanged(const QString &text);

  private:
    void updateResults(qlonglong value);

    QLineEdit *m_binaryInput;
    QLineEdit *m_octalInput;
    QLineEdit *m_decimalInput;
    QLineEdit *m_hexInput;

    bool m_updating;
};

#endif
