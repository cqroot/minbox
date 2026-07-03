#include "mainwindow.h"

#include <QApplication>
#include <QLineEdit>
#include <QListWidget>
#include <QStackedWidget>
#include <QVBoxLayout>
#include <QHBoxLayout>

#include "baseconvertertool.h"
#include "storageconvertertool.h"

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , m_centralWidget(nullptr)
    , m_searchInput(nullptr)
    , m_toolList(nullptr)
    , m_contentStack(nullptr)
    , m_baseConverter(nullptr)
    , m_storageConverter(nullptr)
{
    setupUi();
}

MainWindow::~MainWindow()
{
}

void MainWindow::setupUi()
{
    resize(700, 500);
    setWindowTitle("MinBox - Programmer Toolbox");

    m_centralWidget = new QWidget(this);
    setCentralWidget(m_centralWidget);

    QHBoxLayout *mainLayout = new QHBoxLayout(m_centralWidget);

    QWidget *sidebarWidget = new QWidget();
    sidebarWidget->setMinimumSize(180, 0);
    sidebarWidget->setMaximumSize(220, 16777215);
    QVBoxLayout *sidebarLayout = new QVBoxLayout(sidebarWidget);

    m_searchInput = new QLineEdit();
    m_searchInput->setPlaceholderText("Search tools...");
    sidebarLayout->addWidget(m_searchInput);

    m_toolList = new QListWidget();
    m_toolList->setFrameShape(QFrame::NoFrame);
    sidebarLayout->addWidget(m_toolList);

    mainLayout->addWidget(sidebarWidget);

    m_contentStack = new QStackedWidget();
    m_contentStack->setSizePolicy(QSizePolicy::Minimum, QSizePolicy::Minimum);

    QWidget *contentContainer = new QWidget();
    QVBoxLayout *contentContainerLayout = new QVBoxLayout(contentContainer);
    contentContainerLayout->addStretch();

    QHBoxLayout *centerLayout = new QHBoxLayout();
    centerLayout->addStretch();
    centerLayout->addWidget(m_contentStack);
    centerLayout->addStretch();
    contentContainerLayout->addLayout(centerLayout);

    contentContainerLayout->addStretch();
    mainLayout->addWidget(contentContainer);

    m_baseConverter = new BaseConverterTool();
    m_contentStack->addWidget(m_baseConverter);

    m_storageConverter = new StorageConverterTool();
    m_contentStack->addWidget(m_storageConverter);

    m_allTools << "Base Conversion" << "Storage Unit Conversion";
    for (const QString &tool : m_allTools) {
        m_toolList->addItem(tool);
    }

    connect(m_toolList, &QListWidget::currentRowChanged, this, &MainWindow::onToolListClicked);
    connect(m_searchInput, &QLineEdit::textChanged, this, &MainWindow::onSearchTextChanged);

    m_toolList->setCurrentRow(0);
}

void MainWindow::onToolListClicked(int row)
{
    if (row >= 0) {
        m_contentStack->setCurrentIndex(row);
    }
}

void MainWindow::onSearchTextChanged(const QString &text)
{
    m_toolList->clear();
    for (const QString &tool : m_allTools) {
        if (tool.contains(text, Qt::CaseInsensitive)) {
            m_toolList->addItem(tool);
        }
    }
    if (m_toolList->count() > 0) {
        m_toolList->setCurrentRow(0);
    }
}
