#include "baseconvertertool.h"

#include <QLineEdit>
#include <QLabel>
#include <QGroupBox>
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QFormLayout>
#include <QSpacerItem>
#include <QRegularExpressionValidator>

BaseConverterTool::BaseConverterTool(QWidget *parent)
    : QWidget(parent)
    , m_binaryInput(nullptr)
    , m_octalInput(nullptr)
    , m_decimalInput(nullptr)
    , m_hexInput(nullptr)
    , m_updating(false)
{
    setSizePolicy(QSizePolicy::Minimum, QSizePolicy::Minimum);

    QVBoxLayout *wrapperLayout = new QVBoxLayout(this);
    wrapperLayout->addStretch();
    wrapperLayout->setAlignment(wrapperLayout, Qt::AlignCenter);

    QLabel *title = new QLabel("Base Conversion");
    title->setStyleSheet("font-size: 18px; font-weight: bold;");
    wrapperLayout->addWidget(title);

    QFormLayout *formLayout = new QFormLayout();

    m_binaryInput = new QLineEdit();
    m_binaryInput->setPlaceholderText("Binary");
    m_binaryInput->setValidator(new QRegularExpressionValidator(QRegularExpression("[01]*"), this));
    formLayout->addRow(m_binaryInput);

    m_octalInput = new QLineEdit();
    m_octalInput->setPlaceholderText("Octal");
    m_octalInput->setValidator(new QRegularExpressionValidator(QRegularExpression("[0-7]*"), this));
    formLayout->addRow(m_octalInput);

    m_decimalInput = new QLineEdit();
    m_decimalInput->setPlaceholderText("Decimal");
    m_decimalInput->setValidator(new QRegularExpressionValidator(QRegularExpression("[0-9]*"), this));
    formLayout->addRow(m_decimalInput);

    m_hexInput = new QLineEdit();
    m_hexInput->setPlaceholderText("Hexadecimal");
    m_hexInput->setValidator(new QRegularExpressionValidator(QRegularExpression("(0[xX])?[0-9A-Fa-f]*"), this));
    formLayout->addRow(m_hexInput);

    wrapperLayout->addLayout(formLayout);

    wrapperLayout->addStretch();

    connect(m_binaryInput, &QLineEdit::textChanged, this, &BaseConverterTool::onBinaryTextChanged);
    connect(m_octalInput, &QLineEdit::textChanged, this, &BaseConverterTool::onOctalTextChanged);
    connect(m_decimalInput, &QLineEdit::textChanged, this, &BaseConverterTool::onDecimalTextChanged);
    connect(m_hexInput, &QLineEdit::textChanged, this, &BaseConverterTool::onHexadecimalTextChanged);
}

BaseConverterTool::~BaseConverterTool()
{
}

void BaseConverterTool::onBinaryTextChanged(const QString &text)
{
    if (m_updating) {
        return;
    }
    if (text.isEmpty()) {
        updateResults(0);
        return;
    }
    bool ok = false;
    qlonglong value = text.toLongLong(&ok, 2);
    if (ok) {
        m_updating = true;
        m_octalInput->setText(QString::number(value, 8));
        m_decimalInput->setText(QString::number(value, 10));
        m_hexInput->setText(QString::number(value, 16).toUpper());
        m_updating = false;
    }
}

void BaseConverterTool::onOctalTextChanged(const QString &text)
{
    if (m_updating) {
        return;
    }
    if (text.isEmpty()) {
        updateResults(0);
        return;
    }
    bool ok = false;
    qlonglong value = text.toLongLong(&ok, 8);
    if (ok) {
        m_updating = true;
        m_binaryInput->setText(QString::number(value, 2));
        m_decimalInput->setText(QString::number(value, 10));
        m_hexInput->setText(QString::number(value, 16).toUpper());
        m_updating = false;
    }
}

void BaseConverterTool::onDecimalTextChanged(const QString &text)
{
    if (m_updating) {
        return;
    }
    if (text.isEmpty()) {
        updateResults(0);
        return;
    }
    bool ok = false;
    qlonglong value = text.toLongLong(&ok, 10);
    if (ok) {
        m_updating = true;
        m_binaryInput->setText(QString::number(value, 2));
        m_octalInput->setText(QString::number(value, 8));
        m_hexInput->setText(QString::number(value, 16).toUpper());
        m_updating = false;
    }
}

void BaseConverterTool::onHexadecimalTextChanged(const QString &text)
{
    if (m_updating) {
        return;
    }
    if (text.isEmpty()) {
        updateResults(0);
        return;
    }
    QString input = text;
    if (input.startsWith("0x") || input.startsWith("0X")) {
        input = input.mid(2);
    }
    if (input.isEmpty()) {
        updateResults(0);
        return;
    }
    bool ok = false;
    qlonglong value = input.toLongLong(&ok, 16);
    if (ok) {
        m_updating = true;
        m_binaryInput->setText(QString::number(value, 2));
        m_octalInput->setText(QString::number(value, 8));
        m_decimalInput->setText(QString::number(value, 10));
        m_hexInput->blockSignals(true);
        m_hexInput->setText("0x" + input.toUpper());
        m_hexInput->blockSignals(false);
        m_updating = false;
    }
}

void BaseConverterTool::updateResults(qlonglong value)
{
    m_updating = true;
    m_binaryInput->setText(QString::number(value, 2));
    m_octalInput->setText(QString::number(value, 8));
    m_decimalInput->setText(QString::number(value, 10));
    m_hexInput->blockSignals(true);
    m_hexInput->setText("0x" + QString::number(value, 16).toUpper());
    m_hexInput->blockSignals(false);
    m_updating = false;
}
